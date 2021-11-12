package crypto_test

import (
	"github.com/sisukas/commons/crypto"
	"testing"
)

func TestDecryptAES(t *testing.T) {
	tests := []struct {
		crypt string
		key   string
		msg   string
	}{
		{
			"Y11ck0Qghyg-NG6B35rU8Tp8kpdSUrjOxTmp23bsTeW7enFUxWiBLv3-HNgiR90C",
			"12345678901234561234567890123456",
			"This is the message",
		},
		{
			"mM4yos4mcce1oTmjlDGT6l587c0GFxDk_O4YhY7kOGZ4Aa87ge7O1_JmfdBFsG0Q",
			"12345678901234561234567890123456",
			"This is the messageone",
		},
		{
			"0jWXC3DY9BahsgSIYtNIaOmBLJCea_j_D_wrx1qiibY9hdoixIzQr93zvWPPplgT",
			"adsajkhsajkd676wdhsagtwwq54qw7w6",
			"Another Message here?",
		},
		{
			"SEXU_m2k6A_nZjLHsnhiPartScb-NRb0wffACGsw5iHChw4_JVGqor6lWLvUrRPB",
			"gawvhjk60i5syqntibkibhmq0dluaseu",
			"pKs1idilv1YtLjPVebO8FQ",
		},
	}
	for _, test := range tests {
		msg, err := crypto.DecryptAES(test.crypt, test.key)
		if err != nil {
			t.Fatalf("error decrypting %v ", err)
		}
		if msg != test.msg {
			t.Fatalf("decryption failed crypt %s ", test.crypt)
		}
		//t.Logf("Message %s", msg)
	}

}

func TestDecryptAESToHexString(t *testing.T) {
	tests := []struct {
		crypt  string
		key    string
		b64url string
		hexid  string
	}{
		{
			"RWTICbITRX3VhJ7IZ5lE2K2IO9SFfNob1voTvQ0C4jaf7au3RjD0kPW7UYmk_Pjx",
			"jyxwu1dixpbh5mfz7bbnhubkoh3xxz5g",
			"Rqolrsp60ztf8cgs4H_aRQ",
			"46aa25aeca7ad33b5ff1c82ce07fda45",
		},
		{
			"WYirjRm1AyVm6VbdW4DNvjF5tlLkGBMWJOnG4gPrrAOWOJdQOhZ_5lOW5Rf6BrSb",
			"vqfeakclemmajvm347q3alpel3amupk2",
			"x3mn2ke6aQuOl_tP3ixF5g",
			"c779a7da47ba690b8e97fb4fde2c45e6",
		},
	}
	for _, test := range tests {
		b64, err := crypto.DecryptAES(test.crypt, test.key)
		if err != nil {
			t.Fatalf("error decrypting %v ", err)
		}
		if test.b64url != b64 {
			t.Fatalf("error decrypting  expected %v  got %v", test.b64url, b64)
		}
		ss, err := crypto.URLEncodedBase64ToHex(b64)
		if err != nil {
			t.Fatalf("Error decoding from b64 to hex %v", err)
		}
		if test.hexid != ss {
			t.Fatalf("error decoding to hex  expected %v  got %v", test.hexid, ss)
		}
		//t.Logf("Message %s", msg)
	}

}

func TestEncryptAES(t *testing.T) {

	tests := []struct {
		key string
		msg string
	}{
		{
			"J7IZ5lE2K2IO9SFfNob1voTvQ0C4jafN",
			"This is the message",
		},
		{
			"46aa25aeca7ad33b5ff1c82ce07fda45",
			"Finding another message",
		},
		{
			"SIYtNIaOmBLJCea_j_D_wrx1qiibY9hd",
			"J7IZ5lE2K2IO9SFfNob1voTvQ",
		},
	}

	for _, test := range tests {

		encr, err := crypto.EncryptAES(test.msg, test.key)
		if err != nil {
			t.Fatalf("Error encrypting %v ", err)
		}
		t.Logf("Encrypted message %s", encr)

		msg2, err := crypto.DecryptAES(encr, test.key)
		if err != nil {
			t.Fatalf("Error decrypting %v ", err)
		}
		t.Logf("Decrypted message %s", msg2)

		if test.msg != msg2 {
			t.Fatal("The decryption failed")
		}
	}

}
