package provider

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type KMSProvider struct {
	GenericProvider
	data       kmsProviderData
	region     string
	ciphertext string
	datakey    string
}

type kmsProviderData struct {
	kms.DecryptOutput
}

func NewKMSProvider(ch *chan Provider, params map[string]string) Provider {
	return &KMSProvider{
		region:     params["region"],
		ciphertext: params["ciphertext"],
		datakey:    params["datakey"],
		GenericProvider: GenericProvider{
			renewable: false,
			channel:   ch,
		},
	}
}

func (k *KMSProvider) Initialize() error {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(k.region)}))
	svc := kms.New(sess)

	input := &kms.DecryptInput{}
	if k.datakey != "" {
		input.SetCiphertextBlob([]byte(k.datakey))
	} else {
		input.SetCiphertextBlob([]byte(k.ciphertext))
	}

	data, err := svc.Decrypt(input)
	if err != nil {
		return err
	}

	if k.datakey == "" {
		k.data.Plaintext = data.Plaintext
		k.data.KeyId = data.KeyId
		k.data.DecryptOutput = *data
		return nil
	}

	//const decipher = crypto.createDecipher('aes-256-cbc', data.data.Plaintext);
	block, err := aes.NewCipher(data.Plaintext)
	if err != nil {
		return err
	}

	//let decryptedPlaintext = decipher.update(this.payload.ciphertext, 'base64', 'base64');
	//decryptedPlaintext += decipher.final('base64');
	ciphertext := []byte(k.ciphertext)

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	k.data.Plaintext = ciphertext
	k.data.KeyId = data.KeyId
	return nil
}

func (k *KMSProvider) Renew() error {
	return k.Initialize()
}

func (k *KMSProvider) Invalidate() {}

func (k *KMSProvider) Get() (map[string]string, error) {
	return map[string]string{
		"plaintext": string(k.data.Plaintext),
	}, nil
}
