package storage

func GetDescription(transactions []Transaction) []string {
	descriptions := make([]string, 0, len(transactions))
	for _, transaction := range transactions {
		descriptions = append(descriptions, transaction.Description)
	}

	return descriptions
}
