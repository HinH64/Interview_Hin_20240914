package enums

type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodThirdParty   PaymentMethod = "third_party"
	PaymentMethodBlockchain   PaymentMethod = "blockchain"
)

func (p PaymentMethod) IsValid() bool {
	switch p {
	case PaymentMethodCreditCard, PaymentMethodBankTransfer, PaymentMethodThirdParty, PaymentMethodBlockchain:
		return true
	}
	return false
}