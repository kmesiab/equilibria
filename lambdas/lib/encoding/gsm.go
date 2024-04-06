package encoding

// Standard GSM 03.38 characters (including basic punctuation and control characters)
const standardGSMChars = "@£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞÆæßÉ !\"#¤%&'()*+,-./:;<=>?¡ÄÖÑÜ§¿äöñüà0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Extended GSM 03.38 characters that require an escape character
const extendedGSMChars = "{}\\[~]|€^"

// isStandardGSMChar checks if a rune is in the standard GSM 03.38 character set.
func isStandardGSMChar(r rune) bool {
	for _, c := range standardGSMChars {
		if r == c {
			return true
		}
	}
	return false
}

// isExtendedGSMChar checks if a rune is in the extended GSM 03.38 character set.
func isExtendedGSMChar(r rune) bool {
	for _, c := range extendedGSMChars {
		if r == c {
			return true
		}
	}
	return false
}

// IsGSMEncoded checks if a text is fully GSM 03.38 encoded, considering both standard and extended characters.
func IsGSMEncoded(text string) bool {
	for _, r := range text {
		if !isStandardGSMChar(r) && !isExtendedGSMChar(r) {
			return false
		}
	}
	return true
}
