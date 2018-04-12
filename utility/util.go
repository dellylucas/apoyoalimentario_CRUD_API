package utility

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

//Semester - calcula el semestre actual
//Param round		OUT   "semestre actual"
func Semester() (round int) {
	NowMonth := time.Now().UTC().Month()

	if NowMonth == 1 || NowMonth == 2 || NowMonth == 3 || NowMonth == 4 || NowMonth == 5 || NowMonth == 6 {
		round = 1
	} else {
		round = 3
	}
	return round
}

//GetServiceXML - Obtiene datos desde una URL de un servicio XML
//Param T		IN   "modelo que va a ser el mmapeo de los datos retornados"
//Param url		IN   "url a realizar la peticion"
//Param wg	IN   "variable para alertar las gorutines"
//Param err		OUT   "error si es que existe"
func GetServiceXML(T interface{}, url string, wg *sync.WaitGroup) error {

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("GET error: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Status error: %v", response.StatusCode)
	}
	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("Read body: %v", err)
	}
	xml.Unmarshal(data, &T)
	if wg != nil {
		wg.Done()
	}
	return err
}

//GetInitEnd - Retorna el comienzo y fin de un semestre en formato Time
//Param fromDate		OUT   "fecha comienzo de semestre actual"
//Param toDate		OUT   "fecha fin de semestre actual"
func GetInitEnd() (fromDate time.Time, toDate time.Time) {
	Semester := Semester()
	var Inicial time.Month
	var Final time.Month
	if Semester == 1 {
		Inicial = time.January
		Final = time.June
	} else {
		Inicial = time.July
		Final = time.December
	}
	fromDate = time.Date(time.Now().UTC().Year(), Inicial, 1, 0, 0, 0, 0, time.UTC)
	toDate = time.Date(time.Now().UTC().Year(), Final, 30, 0, 0, 0, 0, time.UTC)
	return fromDate, toDate
}

//createHash - crea hash es la llave para desencriptacion
//param key in "llave de desencriptacion"
//param  out "hash"
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

//Encrypt - encripta contraseña correo
//param data in "dato a encriptar"
//param passphrase in "frase llave de desencriptacion"
//param ciphertext out "encriptacion"
func Encrypt(data []byte, passphrase string) string {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return fmt.Sprintf("%s", ciphertext)
}

//Decrypt - desencripta contraseña correo
//param data in "encriptacion"
//param passphrase in "frase llave de desencriptacion"
//param plaintext out "desencriptacion"
func Decrypt(data []byte, passphrase string) string {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%s", plaintext)
}
