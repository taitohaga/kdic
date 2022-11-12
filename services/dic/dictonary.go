package dic

type dictionaryService struct {}

type GetDictionaryResponse struct {
    dic_name string
    author_name string
    word_count uint64
}

type DictionaryService interface {
    GetDictionary(dic_name string) GetDictionaryResponse
}
