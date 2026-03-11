package dto

type CreateBookRequest struct {
	Title       string `json:"title" validate:"required" example:"Harry Potter and the Philosopher's Stone"`
	Description string `json:"description" validate:"required" example:"Harry Potter has never even heard of Hogwarts when the letters start dropping on the doormat at number four, Privet Drive. Addressed in green ink on yellowish parchment with a purple seal, they are swiftly confiscated by his grisly aunt and uncle."`
	AuthorID    string `json:"author_id" validate:"required" example:"019cd34a-e108-7c85-936e-fee8037b391c"`
}

type UpdateBookRequest struct {
	Title       string `json:"title" validate:"required" example:"Harry Potter and the Philosopher's Stone"`
	Description string `json:"description" validate:"required" example:"Harry Potter has never even heard of Hogwarts when the letters start dropping on the doormat at number four, Privet Drive. Addressed in green ink on yellowish parchment with a purple seal, they are swiftly confiscated by his grisly aunt and uncle."`
}
