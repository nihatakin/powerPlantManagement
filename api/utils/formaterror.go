package formaterror

import (
	"errors"
	"strings"
)

//TODO eksik kaldı burası doldurmak gerekiyor.
func FormatError(err string) error {
	if strings.Contains(err, "name") {
		return errors.New("Name Already Taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}
	if strings.Contains(err, "shortName") {
		return errors.New("ShortName Already Taken")
	}
	if strings.Contains(err, "etsoCode") {
		return errors.New("EtsoCode Already Taken")
	}
	if strings.Contains(err, "installedPower") {
		return errors.New("installedPower is required")
	}
	if strings.Contains(err, "order") {
		return errors.New("order is required")
	}
	if strings.Contains(err, "password") {
		return errors.New("Incorrect Password")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}
	return errors.New("Incorrect Details")
}