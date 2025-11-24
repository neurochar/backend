package traits

import "github.com/neurochar/backend/internal/domain/testing/entity"

var Traits = []entity.PersonalityTrait{
	&entity.PersonalityTraitBipolar{
		ID:   10,
		Name: "Теплота",
		// nolint
		Description:    `Отражает степень ориентированности на людей. Высокая теплота — это не «наивность», а готовность к эмоциональному контакту и сопричастности.`,
		LeftStateName:  "Холодный",
		RightStateName: "Теплый",
	},
	&entity.PersonalityTraitBipolar{
		ID:   11,
		Name: "Интеллект",
		// nolint
		Description:    `Способность оперировать не только реальными предметами, но и умственными моделями. Абстрактность здесь означает умение удерживать несколько элементов в памяти, связывать их, обобщать и мыслить символами, а не то, что человек «витает в облаках».`,
		LeftStateName:  "Конкретный",
		RightStateName: "Абстрактный",
	},
	&entity.PersonalityTraitBipolar{
		ID:   12,
		Name: "Эмоциональная устойчивость",
		// nolint
		Description:    `Способность сохранять внутреннее равновесие. Низкая устойчивость — не «истеричность», а повышенная реактивность к стрессу.`,
		LeftStateName:  "Неустойчивый",
		RightStateName: "Устойчивый",
	},
	&entity.PersonalityTraitBipolar{
		ID:   13,
		Name: "Доминирование",
		// nolint
		Description:    `Уровень напора в социальных взаимодействиях. Высокое доминирование — не агрессия, а уверенное влияние`,
		LeftStateName:  "Уступчивый",
		RightStateName: "Властный",
	},
	&entity.PersonalityTraitBipolar{
		ID:   14,
		Name: "Живость",
		// nolint
		Description:    `Выраженность внешней эмоциональности. Низкая живость — не «тусклость», а сдержанность.`,
		LeftStateName:  "Сдержанный",
		RightStateName: "Живой",
	},
	&entity.PersonalityTraitBipolar{
		ID:   15,
		Name: "Нормативность",
		// nolint
		Description:    `Отношение к правилам и обязанностям. Высокая нормативность — не «заорганизованность», а надёжность и ответственность.`,
		LeftStateName:  "Ненадёжный",
		RightStateName: "Дисциплинированный",
	},
	&entity.PersonalityTraitBipolar{
		ID:   16,
		Name: "Социальная смелость",
		// nolint
		Description:    `Комфортность в ситуациях общения. Низкий уровень — не «социальный страх», а избегание внимания.`,
		LeftStateName:  "Застенчивый",
		RightStateName: "Смелый",
	},
	&entity.PersonalityTraitBipolar{
		ID:   17,
		Name: "Чувствительность",
		// nolint
		Description:    `Мягкость и эмоциональная восприимчивость. Высокая чувствительность — не слабость, а развитая эмпатия.`,
		LeftStateName:  "Жесткий",
		RightStateName: "Чувствительный",
	},
	&entity.PersonalityTraitBipolar{
		ID:   18,
		Name: "Подозрительность",
		// nolint
		Description:    `Степень доверия к людям. Высокий уровень — не «паранойя», а настороженная критичность.`,
		LeftStateName:  "Доверчивый",
		RightStateName: "Подозрительный",
	},
	&entity.PersonalityTraitBipolar{
		ID:   19,
		Name: "Мечтательность",
		// nolint
		Description:    `Склонность жить идеями. Высокий уровень — не «оторванность», а развитое внутреннее воображение.`,
		LeftStateName:  "Практичный",
		RightStateName: "Мечтательный",
	},
	&entity.PersonalityTraitBipolar{
		ID:   20,
		Name: "Проницательность",
		// nolint
		Description:    `Социальная наблюдательность и понимание намерений людей. Высокий уровень — не «манипулятивность», а дипломатичность.`,
		LeftStateName:  "Прямой",
		RightStateName: "Проницательный",
	},
	&entity.PersonalityTraitBipolar{
		ID:   21,
		Name: "Тревожность",
		// nolint
		Description:    `Внутренняя напряжённость и склонность к самокритике. Низкий уровень — не «безразличие», а спокойствие.`,
		LeftStateName:  "Уверенный",
		RightStateName: "Тревожный",
	},
	&entity.PersonalityTraitBipolar{
		ID:   22,
		Name: "Открытость изменениям",
		// nolint
		Description:    `Отношение к новому. Высокий уровень — не «бунтарство», а готовность пересматривать взгляды.`,
		LeftStateName:  "Консервативный",
		RightStateName: "Новаторский",
	},
	&entity.PersonalityTraitBipolar{
		ID:   23,
		Name: "Самостоятельность",
		// nolint
		Description:    `Степень опоры на собственные решения. Высокий уровень — не «индивидуализм», а автономия.`,
		LeftStateName:  "Зависимый",
		RightStateName: "Независимый",
	},
	&entity.PersonalityTraitBipolar{
		ID:   24,
		Name: "Самоконтроль",
		// nolint
		Description:    `Организованность и способность регулировать поведение. Высокий уровень — не «зажатость», а порядок и самодисциплина.`,
		LeftStateName:  "Несобранный",
		RightStateName: "Организованный",
	},
	&entity.PersonalityTraitBipolar{
		ID:   25,
		Name: "Напряженность",
		// nolint
		Description:    `Внутреннее физиологическое и психологическое возбуждение. Высокий уровень — не «паника», а высокий тонус и взвинченность.`,
		LeftStateName:  "Спокойный",
		RightStateName: "Напряженный",
	},
}
