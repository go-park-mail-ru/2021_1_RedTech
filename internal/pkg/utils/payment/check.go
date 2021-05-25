package payment

import (
	"Redioteka/internal/pkg/subscription/delivery/grpc/proto"
	"Redioteka/internal/pkg/utils/log"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"strings"
)

func Check(form *proto.Payment) error {
	if form.Unaccepted && !form.CodePro {
		return errors.New("Payment is unaccepted")
	} else if form.CodePro {
		log.Log.Info("Payment from wallet is protected - need to input code")
	}
	if form.Label == "" {
		return errors.New("Label of payment is empty")
	}
	return CheckHash(form)
}

func CheckHash(form *proto.Payment) error {
	secret, err := ioutil.ReadFile("secret")
	if err != nil {
		return err
	}

	codepro := "true"
	if !form.CodePro {
		codepro = "false"
	}
	params := strings.Join([]string{form.Type, form.OperationID, form.Amount, form.Currency,
		form.DateTime, form.Sender, codepro, string(secret), form.Label}, "&")
	hash := sha1.Sum([]byte(params))
	hexString := hex.EncodeToString(hash[:])
	if hexString != form.Hash {
		return errors.New("Payment hash do not match")
	}
	return nil
}
