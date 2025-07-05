package mock

type GetDoctorsRequest struct {
	Stomatology int `json:"stomatology" validate:"required"`
}

type Doctor struct {
	ID          string   `json:"id"`
	Stomatology int      `json:"stomatology"`
	FullName    string   `json:"full_name"`
	Specialties []string `json:"specialties"`
	Filial      string   `json:"filial"`
	Description string   `json:"description"`
}

func (s *Service) GetDoctors(request *GetDoctorsRequest) []Doctor {
	doctors := []Doctor{
		{
			ID:          "1",
			Stomatology: 2,
			FullName:    "Бекишев Серик Асылханович",
			Specialties: []string{"Ортодонтия", "Терапевтическая стоматология"},
			Filial:      "Мангилик Ел 55/15",
			Description: "Опытный стоматолог с уклоном в ортодонтию и лечение кариеса.",
		},
		{
			ID:          "2",
			Stomatology: 2,
			FullName:    "Данияр Оразалиевич",
			Specialties: []string{"Детская стоматология", "Профилактика"},
			Filial:      "Мангилик Ел 55/15",
			Description: "Специалист по детской стоматологии и профилактике заболеваний полости рта.",
		},
		{
			ID:          "3",
			Stomatology: 2,
			FullName:    "Абинаева Жаннат Смагуловна",
			Specialties: []string{"Хирургическая стоматология", "Имплантология"},
			Filial:      "Мангилик Ел 55/15",
			Description: "Проводит сложные хирургические вмешательства и установку имплантов.",
		},
		{
			ID:          "4",
			Stomatology: 2,
			FullName:    "Садырова Асель",
			Specialties: []string{"Эстетическая стоматология", "Отбеливание"},
			Filial:      "Бухар Жырау 19",
			Description: "Специализируется на эстетической стоматологии и отбеливании зубов.",
		},
		{
			ID:          "5",
			Stomatology: 2,
			FullName:    "Тусупов Ерлан",
			Specialties: []string{"Пародонтология", "Гигиена полости рта"},
			Filial:      "Бухар Жырау 19",
			Description: "Занимается лечением заболеваний десен и гигиеной полости рта.",
		},
	}

	var filteredDoctors []Doctor
	for _, doctor := range doctors {
		if doctor.Stomatology == request.Stomatology {
			filteredDoctors = append(filteredDoctors, doctor)
		}
	}

	return filteredDoctors
}
