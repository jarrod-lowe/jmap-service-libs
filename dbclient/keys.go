package dbclient

// Primary key attribute names.
const (
	AttrPK = "pk"
	AttrSK = "sk"
)

// Common key prefixes shared across all JMAP services.
const (
	PrefixAccount = "ACCOUNT#"
	PrefixUser    = "USER#"
	SKMeta        = "META#"
)

// AccountPK returns the partition key for an account.
func AccountPK(accountID string) string {
	return PrefixAccount + accountID
}

// UserPK returns the partition key for a user.
func UserPK(userID string) string {
	return PrefixUser + userID
}
