package kettel

var ItemsLib = map[uint64]*KettelItemAnswerImpl{}

func init() {
	ItemsLib[3] = &KettelItemAnswerImpl{
		ID:      3,
		TraitID: 10,
		// nolint
		MaleQuestion: `Я бы предпочел временами жить в месте, которое находится:`,
		// nolint
		FemaleQuestion: `Я бы предпочла временами жить в месте, которое находится:`,
		Variants: []string{
			`в мегаполисе`,
			`нечто среднее`,
			`в спокойном месте вдали от городской суеты`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[4] = &KettelItemAnswerImpl{
		ID:      4,
		TraitID: 12,
		// nolint
		MaleQuestion: `Я чувствую в себе достаточно сил, чтобы справиться со своими трудностями:`,
		Variants: []string{
			`всегда`,
			`обычно`,
			`редко`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[5] = &KettelItemAnswerImpl{
		ID:      5,
		TraitID: 12,
		// nolint
		MaleQuestion: `Я чувствую некоторое беспокойство при виде диких животных, даже если они находятся в прочных клетках:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[6] = &KettelItemAnswerImpl{
		ID:      6,
		TraitID: 13,
		// nolint
		MaleQuestion: `Я воздерживаюсь от критики людей и их высказываний:`,
		Variants: []string{
			`да`,
			`иногда`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[7] = &KettelItemAnswerImpl{
		ID:      7,
		TraitID: 13,
		// nolint
		MaleQuestion: `Я делаю саркастические (язвительные) замечания по поводу людей, если они этого, по-моему, заслуживают:`,
		Variants: []string{
			`обычно`,
			`иногда`,
			`никогда`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[8] = &KettelItemAnswerImpl{
		ID:      8,
		TraitID: 14,
		// nolint
		MaleQuestion: `Мне больше нравится классическая, чем популярная музыка:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[9] = &KettelItemAnswerImpl{
		ID:      9,
		TraitID: 15,
		// nolint
		MaleQuestion: `Если бы я увидел, как дерутся соседские дети, то я:`,
		// nolint
		FemaleQuestion: `Если бы я увидела, как дерутся соседские дети, то я:`,
		Variants: []string{
			`дал бы им возможность разобраться самим`,
			`не уверен`,
			`попытался бы их успокоить`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[10] = &KettelItemAnswerImpl{
		ID:      10,
		TraitID: 16,
		// nolint
		MaleQuestion: `При общении с людьми я:`,
		Variants: []string{
			`с готовностью вступаю в разговор`,
			`нечто среднее`,
			`предпочитаю спокойно оставаться в стороне`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[11] = &KettelItemAnswerImpl{
		ID:      11,
		TraitID: 17,
		// nolint
		MaleQuestion: `По-моему, интереснее быть:`,
		Variants: []string{
			`разработчиком`,
			`не уверен`,
			`контент-креатором`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[12] = &KettelItemAnswerImpl{
		ID:      12,
		TraitID: 17,
		// nolint
		MaleQuestion: `Я остановился бы на улице скорее, чтобы посмотреть на работу художника, чем слушать, как ссорятся люди:`,
		// nolint
		FemaleQuestion: `Я остановилась бы на улице скорее, чтобы посмотреть на работу художника, чем слушать, как ссорятся люди:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[13] = &KettelItemAnswerImpl{
		ID:      13,
		TraitID: 18,
		// nolint
		MaleQuestion: `Обычно я могу ладить с самодовольными людьми, несмотря на то, что они хвастаются или слишком много о себе воображают:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[14] = &KettelItemAnswerImpl{
		ID:      14,
		TraitID: 19,
		// nolint
		MaleQuestion: `По лицу человека всегда можно заметить, что он нечестный:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[15] = &KettelItemAnswerImpl{
		ID:      15,
		TraitID: 19,
		// nolint
		MaleQuestion: `Для меня важно регулярно делать паузы и восстанавливать силы:`,
		Variants: []string{
			`согласен`,
			`не уверен`,
			`не согласен`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[16] = &KettelItemAnswerImpl{
		ID:      16,
		TraitID: 20,
		// nolint
		MaleQuestion: `Я предпочел бы работу с возможно большим, но непостоянным заработком, чем работу со скромным, но постоянным окладом:`,
		// nolint
		FemaleQuestion: `Я предпочла бы работу с возможно большим, но непостоянным заработком, чем работу со скромным, но постоянным окладом:`,
		Variants: []string{
			`согласен`,
			`не уверен`,
			`не согласен`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[17] = &KettelItemAnswerImpl{
		ID:      17,
		TraitID: 20,
		// nolint
		MaleQuestion: `Я открыто говорю о своих чувствах:`,
		Variants: []string{
			`только если это необходимо`,
			`нечто среднее`,
			`охотно, когда представляется возможность`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[18] = &KettelItemAnswerImpl{
		ID:      18,
		TraitID: 21,
		// nolint
		MaleQuestion: `Я часто мысленно готовлюсь к тому, что что-то может пойти не так:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[19] = &KettelItemAnswerImpl{
		ID:      19,
		TraitID: 21,
		// nolint
		MaleQuestion: `Когда меня неправильно критикуют за что-то, в чем я не виноват, я:`,
		// nolint
		FemaleQuestion: `Когда меня неправильно критикуют за что-то, в чем я не виновата, я:`,
		Variants: []string{
			`не испытываю чувства вины`,
			`нечто среднее`,
			`все же чувствую себя немного виноватым`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[20] = &KettelItemAnswerImpl{
		ID:      20,
		TraitID: 22,
		// nolint
		MaleQuestion: `За деньги можно купить почти все:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[21] = &KettelItemAnswerImpl{
		ID:      21,
		TraitID: 22,
		// nolint
		MaleQuestion: `При принятии важных решений я чаще опираюсь на:`,
		Variants: []string{
			`свои ощущения и интуицию`,
			`и интуицию, и анализ`,
			`факты и логический анализ`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[22] = &KettelItemAnswerImpl{
		ID:      22,
		TraitID: 23,
		// nolint
		MaleQuestion: `Большинство людей были бы больше счастливы, если бы они были ближе друг к другу и поступали так же, как все:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[23] = &KettelItemAnswerImpl{
		ID:      23,
		TraitID: 24,
		// nolint
		MaleQuestion: `Иногда я ловлю себя на том, что делаю что-то «на автомате»:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[24] = &KettelItemAnswerImpl{
		ID:      24,
		TraitID: 24,
		// nolint
		MaleQuestion: `Когда я что-то объясняю, я обычно:`,
		Variants: []string{
			`говорю по ходу мысли`,
			`зависит от ситуации`,
			`стараюсь заранее выстроить план объяснения`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[25] = &KettelItemAnswerImpl{
		ID:      25,
		TraitID: 25,
		// nolint
		MaleQuestion: `Некоторые разговоры или ситуации могут ещё какое-то время оставаться у меня в голове:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[26] = &KettelItemAnswerImpl{
		ID:      26,
		TraitID: 10,
		// nolint
		MaleQuestion: `При одинаковом рабочем времени и заработке было бы интереснее работать:`,
		Variants: []string{
			`программистом`,
			`не уверен`,
			`креативным директором`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[27] = &KettelItemAnswerImpl{
		ID:      27,
		TraitID: 10,
		// nolint
		MaleQuestion: `В группе меня обычно выбирали для организации мероприятий:`,
		Variants: []string{
			`очень редко`,
			`иногда`,
			`много раз`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[28] = &KettelItemAnswerImpl{
		ID:      28,
		TraitID: 11,
		// nolint
		MaleQuestion: `«Браузер» относится к «поиску информации», как «редактор кода» относится к:`,
		Variants: []string{
			`«написанию программ»`,
			`«рисованию»`,
			`«проверке орфографии»`,
		},
		RawVariantKeys: []int{0, 1, 0},
	}

	ItemsLib[29] = &KettelItemAnswerImpl{
		ID:      29,
		TraitID: 12,
		// nolint
		MaleQuestion: `Иногда я не могу заснуть потому что какая-нибудь мысль не выходит из головы:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[30] = &KettelItemAnswerImpl{
		ID:      30,
		TraitID: 12,
		// nolint
		MaleQuestion: `Если цель для меня действительно важна, я обычно продолжаю двигаться к ней несмотря на трудности:`,
		Variants: []string{
			`чаще продолжаю искать способ добиться результата`,
			`зависит от обстоятельств`,
			`иногда считаю разумнее переключиться на что-то другое`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[31] = &KettelItemAnswerImpl{
		ID:      31,
		TraitID: 13,
		// nolint
		MaleQuestion: `Устаревший закон следует изменить:`,
		Variants: []string{
			`только после основательного обсуждения`,
			`не уверен`,
			`как можно скорее`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[32] = &KettelItemAnswerImpl{
		ID:      32,
		TraitID: 13,
		// nolint
		MaleQuestion: `Я чувствую себя неловко, когда мне приходится работать над чем-то, что требует быстрых действий и может повлиять на других:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[33] = &KettelItemAnswerImpl{
		ID:      33,
		TraitID: 14,
		// nolint
		MaleQuestion: `Большинство знакомых считают меня интересным рассказчиком:`,
		// nolint
		FemaleQuestion: `Большинство знакомых считают меня интересной рассказчицей:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[34] = &KettelItemAnswerImpl{
		ID:      34,
		TraitID: 15,
		// nolint
		MaleQuestion: `Когда люди выглядят неаккуратно или не следят за собой, я обычно:`,
		Variants: []string{
			`не придаю этому большого значения`,
			`отношусь по-разному в зависимости от ситуации`,
			`замечаю это и внутренне обращаю на это внимание`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[35] = &KettelItemAnswerImpl{
		ID:      35,
		TraitID: 16,
		// nolint
		MaleQuestion: `Я чувствую себя немного не по себе, если неожиданно оказываюсь в центре внимания группы людей:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[36] = &KettelItemAnswerImpl{
		ID:      36,
		TraitID: 16,
		// nolint
		MaleQuestion: `Я всегда рад оказаться среди людей, например, на вечеринке, концерте, совместном мероприятии:`,
		// nolint
		FemaleQuestion: `Я всегда рада оказаться среди людей, например, на вечеринке, концерте, совместном мероприятии:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[37] = &KettelItemAnswerImpl{
		ID:      37,
		TraitID: 17,
		// nolint
		MaleQuestion: `В школе/университете мне было интереснее:`,
		Variants: []string{
			`придумывать и выражать идеи`,
			`и то, и другое`,
			`разбираться, как всё устроено на практике`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[38] = &KettelItemAnswerImpl{
		ID:      38,
		TraitID: 18,
		// nolint
		MaleQuestion: `Если меня назначают руководителем чего-либо, я настаиваю на том, чтобы мои указания выполнялись, иначе я отказываюсь от этой работы:`,
		Variants: []string{
			`да`,
			`иногда`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[39] = &KettelItemAnswerImpl{
		ID:      39,
		TraitID: 19,
		// nolint
		MaleQuestion: `Важнее, чтобы родители:`,
		Variants: []string{
			`помогали детям развивать свои чувства`,
			`нечто среднее`,
			`обучали детей сдерживать свои чувства`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[40] = &KettelItemAnswerImpl{
		ID:      40,
		TraitID: 19,
		// nolint
		MaleQuestion: `Участвуя в групповой деятельности, я бы предпочел:`,
		// nolint
		FemaleQuestion: `Участвуя в групповой деятельности, я бы предпочла:`,
		Variants: []string{
			`постараться улучшить организацию работы`,
			`нечто среднее`,
			`следить за результатами и соблюдением правил`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[41] = &KettelItemAnswerImpl{
		ID:      41,
		TraitID: 20,
		// nolint
		MaleQuestion: `Мне сложно долгое время чувствовать себя комфортно без движения и активности:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[42] = &KettelItemAnswerImpl{
		ID:      42,
		TraitID: 20,
		// nolint
		MaleQuestion: `В общении мне обычно комфортнее с людьми, которые:`,
		Variants: []string{
			`стараются избегать острых споров`,
			`бывают разными`,
			`любят активно спорить и отстаивать свою точку зрения`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[43] = &KettelItemAnswerImpl{
		ID:      43,
		TraitID: 21,
		// nolint
		MaleQuestion: `Я чувствую себя очень униженным, когда меня критикуют в присутствии группы людей:`,
		// nolint
		FemaleQuestion: `Я чувствую себя очень униженной, когда меня критикуют в присутствии группы людей:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[44] = &KettelItemAnswerImpl{
		ID:      44,
		TraitID: 21,
		// nolint
		MaleQuestion: `Если мне внезапно приходит сообщение «нужно обсудить один вопрос», я обычно:`,
		Variants: []string{
			`не придаю этому особого значения`,
			`зависит от ситуации`,
			`начинаю заранее обдумывать, что могло случиться`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[45] = &KettelItemAnswerImpl{
		ID:      45,
		TraitID: 22,
		// nolint
		MaleQuestion: `В наше время требуется:`,
		Variants: []string{
			`больше спокойных, стабильных людей`,
			`не уверен`,
			`больше «идеалистов», планирующих лучшее будущее`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[46] = &KettelItemAnswerImpl{
		ID:      46,
		TraitID: 22,
		// nolint
		MaleQuestion: `При чтении я сразу замечаю, когда автор произведения хочет меня в чем-то убедить:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[47] = &KettelItemAnswerImpl{
		ID:      47,
		TraitID: 23,
		// nolint
		MaleQuestion: `Мне обычно было интересно участвовать в командных или соревновательных активностях:`,
		Variants: []string{
			`не очень`,
			`по-разному`,
			`да, это меня увлекает`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[48] = &KettelItemAnswerImpl{
		ID:      48,
		TraitID: 24,
		// nolint
		MaleQuestion: `Я поддерживаю порядок в моей комнате, все вещи всегда лежат на своих местах:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[49] = &KettelItemAnswerImpl{
		ID:      49,
		TraitID: 25,
		// nolint
		MaleQuestion: `Иногда у меня возникает чувство напряжения и беспокойства, когда я вспоминаю, что произошло в течение дня:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[50] = &KettelItemAnswerImpl{
		ID:      50,
		TraitID: 25,
		// nolint
		MaleQuestion: `Иногда я сомневаюсь, действительно ли люди, с которыми я разговариваю, интересуются тем, что я говорю:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[51] = &KettelItemAnswerImpl{
		ID:      51,
		TraitID: 10,
		// nolint
		MaleQuestion: `Если бы пришлось выбирать, то я предпочел бы быть:`,
		// nolint
		FemaleQuestion: `Если бы пришлось выбирать, то я предпочла бы быть:`,
		Variants: []string{
			`лидером общественного движения`,
			`не уверен`,
			`экспертом в своей профессиональной области`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[52] = &KettelItemAnswerImpl{
		ID:      52,
		TraitID: 10,
		// nolint
		MaleQuestion: `Подготовка подарков для других людей для меня чаще:`,
		Variants: []string{
			`приятная часть праздника`,
			`по-разному`,
			`скорее формальность`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[53] = &KettelItemAnswerImpl{
		ID:      53,
		TraitID: 11,
		// nolint
		MaleQuestion: `«Выгоревший» относится к «переработкам», как «гордый» относится к:`,
		Variants: []string{
			`«уверенности»`,
			`«достижению»`,
			`«радости»`,
		},
		RawVariantKeys: []int{0, 1, 0},
	}

	ItemsLib[54] = &KettelItemAnswerImpl{
		ID:      54,
		TraitID: 11,
		// nolint
		MaleQuestion: `Какой из следующих предметов по существу отличается от двух других:`,
		Variants: []string{
			`свеча`,
			`луна`,
			`электрический свет`,
		},
		RawVariantKeys: []int{0, 1, 0},
	}

	ItemsLib[55] = &KettelItemAnswerImpl{
		ID:      55,
		TraitID: 12,
		// nolint
		MaleQuestion: `Друзья меня подводили:`,
		Variants: []string{
			`очень редко`,
			`иногда`,
			`довольно часто`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[56] = &KettelItemAnswerImpl{
		ID:      56,
		TraitID: 13,
		// nolint
		MaleQuestion: `У меня есть качества, по которым я определенно выше большинства людей:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[57] = &KettelItemAnswerImpl{
		ID:      57,
		TraitID: 13,
		// nolint
		MaleQuestion: `Когда я расстроен, я стараюсь скрыть свои чувства от других:`,
		// nolint
		FemaleQuestion: `Когда я расстроена, я стараюсь скрыть свои чувства от других:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[58] = &KettelItemAnswerImpl{
		ID:      58,
		TraitID: 14,
		// nolint
		MaleQuestion: `Я склонен посещать мероприятия и развлечения:`,
		// nolint
		FemaleQuestion: `Я склонна посещать мероприятия и развлечения:`,
		Variants: []string{
			`чаще раза в неделю (т.е. чаще, чем большинство)`,
			`примерно раз в неделю (т.е. как большинство)`,
			`реже раза в неделю (т.е. реже, чем большинство)`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[59] = &KettelItemAnswerImpl{
		ID:      59,
		TraitID: 15,
		// nolint
		MaleQuestion: `Я считаю, что возможность вести себя непринужденно важнее, чем хорошие манеры и следование всем правилам:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[60] = &KettelItemAnswerImpl{
		ID:      60,
		TraitID: 16,
		// nolint
		MaleQuestion: `Обычно я молчу в присутствии старших по возрасту, опыту и положению:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[61] = &KettelItemAnswerImpl{
		ID:      61,
		TraitID: 16,
		// nolint
		MaleQuestion: `Мне проще общаться один на один, чем выступать перед большой аудиторией:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[62] = &KettelItemAnswerImpl{
		ID:      62,
		TraitID: 17,
		// nolint
		MaleQuestion: `Если я еду в место, где уже бывал раньше, я чаще:`,
		Variants: []string{
			`еду по памяти`,
			`зависит от ситуации`,
			`на всякий случай открываю навигатор`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[63] = &KettelItemAnswerImpl{
		ID:      63,
		TraitID: 18,
		// nolint
		MaleQuestion: `Если кто-нибудь рассердится на меня, то я:`,
		Variants: []string{
			`постараюсь его успокоить`,
			`нечто среднее`,
			`раздражаюсь`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[64] = &KettelItemAnswerImpl{
		ID:      64,
		TraitID: 18,
		// nolint
		MaleQuestion: `Даже если я с чем-то не согласен, я иногда предпочитаю не вступать в спор без необходимости:`,
		// nolint
		FemaleQuestion: `Даже если я с чем-то не согласена, я иногда предпочитаю не вступать в спор без необходимости:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[65] = &KettelItemAnswerImpl{
		ID:      65,
		TraitID: 19,
		// nolint
		MaleQuestion: `Бывает, что я помню само место или ситуацию, но не могу сразу вспомнить название:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[66] = &KettelItemAnswerImpl{
		ID:      66,
		TraitID: 20,
		// nolint
		MaleQuestion: `Мне бы понравилась работа ветеринара, лечение и операции на животных:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[67] = &KettelItemAnswerImpl{
		ID:      67,
		TraitID: 20,
		// nolint
		MaleQuestion: `В неформальной обстановке я обычно веду себя довольно свободно и непринуждённо:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[68] = &KettelItemAnswerImpl{
		ID:      68,
		TraitID: 21,
		// nolint
		MaleQuestion: `Бывают времена, когда у меня нет настроения видеть кого бы то ни было:`,
		Variants: []string{
			`очень редко`,
			`нечто среднее`,
			`довольно часто`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[69] = &KettelItemAnswerImpl{
		ID:      69,
		TraitID: 21,
		// nolint
		MaleQuestion: `Когда я увлечён разговором или идеей, это обычно заметно со стороны:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[70] = &KettelItemAnswerImpl{
		ID:      70,
		TraitID: 22,
		// nolint
		MaleQuestion: `В юности, если я расходился во мнении с родителями, то я чаще:`,
		// nolint
		FemaleQuestion: `В юности, если я расходилась во мнении с родителями, то я чаще:`,
		Variants: []string{
			`старался(ась) объяснить свою точку зрения`,
			`зависело от ситуации`,
			`предпочитал(а) не затягивать спор`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[71] = &KettelItemAnswerImpl{
		ID:      71,
		TraitID: 23,
		// nolint
		MaleQuestion: `Я предпочел бы заниматься самостоятельной работой, а не совместной с другими:`,
		// nolint
		FemaleQuestion: `Я предпочла бы заниматься самостоятельной работой, а не совместной с другими:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[72] = &KettelItemAnswerImpl{
		ID:      72,
		TraitID: 23,
		// nolint
		MaleQuestion: `Мне бы больше понравилась спокойная жизнь, чем слава и шумный успех:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[73] = &KettelItemAnswerImpl{
		ID:      73,
		TraitID: 24,
		// nolint
		MaleQuestion: `Со временем я стал гораздо спокойнее относиться к большинству жизненных ситуаций:`,
		// nolint
		FemaleQuestion: `Со временем я стала гораздо спокойнее относиться к большинству жизненных ситуаций`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[74] = &KettelItemAnswerImpl{
		ID:      74,
		TraitID: 25,
		// nolint
		MaleQuestion: `Некоторые комментарии или замечания я продолжаю обдумывать уже после завершения разговора:`,
		Variants: []string{
			`часто`,
			`иногда`,
			`никогда`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[75] = &KettelItemAnswerImpl{
		ID:      75,
		TraitID: 25,
		// nolint
		MaleQuestion: `Я всегда способен управлять проявлением своих чувств:`,
		// nolint
		FemaleQuestion: `Я всегда способна управлять проявлением своих чувств:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[76] = &KettelItemAnswerImpl{
		ID:      76,
		TraitID: 10,
		// nolint
		MaleQuestion: `Начиная работу над полезным проектом, я бы предпочел:`,
		// nolint
		FemaleQuestion: `Начиная работу над полезным проектом, я бы предпочла:`,
		Variants: []string{
			`разрабатывать его концепцию`,
			`нечто среднее`,
			`заниматься его практической реализацией`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[77] = &KettelItemAnswerImpl{
		ID:      77,
		TraitID: 11,
		// nolint
		MaleQuestion: `«Удивление» относится к «странный», как «страх» относится к:`,
		Variants: []string{
			`«смелый»`,
			`«тревожный»`,
			`«ужасный»`,
		},
		RawVariantKeys: []int{0, 0, 1},
	}

	ItemsLib[78] = &KettelItemAnswerImpl{
		ID:      78,
		TraitID: 11,
		// nolint
		MaleQuestion: `Которая из последующих дробей отличается от двух других:`,
		Variants: []string{
			`3/7`,
			`3/9`,
			`3/11`,
		},
		RawVariantKeys: []int{0, 1, 0},
	}

	ItemsLib[79] = &KettelItemAnswerImpl{
		ID:      79,
		TraitID: 12,
		// nolint
		MaleQuestion: `Бывает, что в общении с людьми я ощущаю некоторую дистанцию, не всегда понимая её причину:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[80] = &KettelItemAnswerImpl{
		ID:      80,
		TraitID: 12,
		// nolint
		MaleQuestion: `Бывает, что люди понимают мои намерения не совсем так, как мне хотелось бы:`,
		Variants: []string{
			`часто`,
			`иногда`,
			`никогда`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[81] = &KettelItemAnswerImpl{
		ID:      81,
		TraitID: 13,
		// nolint
		MaleQuestion: `В общении я чаще предпочитаю людей, которые стараются обходиться без грубых выражений:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[82] = &KettelItemAnswerImpl{
		ID:      82,
		TraitID: 14,
		// nolint
		MaleQuestion: `У меня определенно меньше друзей, чем у большинства людей:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[83] = &KettelItemAnswerImpl{
		ID:      83,
		TraitID: 14,
		// nolint
		MaleQuestion: `Я бы очень не хотел находиться в таком месте, где нет таких людей, с которыми можно поговорить:`,
		// nolint
		FemaleQuestion: `Я бы очень не хотела находиться в таком месте, где нет таких людей, с которыми можно поговорить:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[84] = &KettelItemAnswerImpl{
		ID:      84,
		TraitID: 15,
		// nolint
		MaleQuestion: `Люди иногда считают меня небрежным, хотя и думают, что я приятный человек:`,
		// nolint
		FemaleQuestion: `Люди иногда считают меня небрежной, хотя и думают, что я приятный человек:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[85] = &KettelItemAnswerImpl{
		ID:      85,
		TraitID: 16,
		// nolint
		MaleQuestion: `Когда нужно выступать перед большой аудиторией, мне обычно требуется время, чтобы внутренне настроиться:`,
		Variants: []string{
			`довольно часто`,
			`иногда`,
			`почти никогда`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[86] = &KettelItemAnswerImpl{
		ID:      86,
		TraitID: 16,
		// nolint
		MaleQuestion: `В большой компании я чаще сначала наблюдаю за разговором, чем сразу активно включаюсь в него:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[87] = &KettelItemAnswerImpl{
		ID:      87,
		TraitID: 17,
		// nolint
		MaleQuestion: `Я скорее предпочту почитать или посмотреть:`,
		Variants: []string{
			`документальные фильмы о науке и технологиях`,
			`нечто среднее`,
			`фильмы и сериалы с глубокими эмоциями`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[88] = &KettelItemAnswerImpl{
		ID:      88,
		TraitID: 18,
		// nolint
		MaleQuestion: `Мне комфортнее, когда со мной договариваются, а не просто говорят, что делать:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[89] = &KettelItemAnswerImpl{
		ID:      89,
		TraitID: 18,
		// nolint
		MaleQuestion: `Когда люди выражают недовольство чем-то в моих действиях, я обычно понимаю причину такой реакции:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[90] = &KettelItemAnswerImpl{
		ID:      90,
		TraitID: 19,
		// nolint
		MaleQuestion: `На улицах или в магазинах мне не нравится, когда некоторые люди пристально разглядывают других:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[91] = &KettelItemAnswerImpl{
		ID:      91,
		TraitID: 19,
		// nolint
		MaleQuestion: `Во время длительной поездки я бы предпочел:`,
		// nolint
		FemaleQuestion: `Во время длительной поездки я бы предпочла:`,
		Variants: []string{
			`читать что-нибудь серьёзное, но интересное`,
			`неопределенно`,
			`провести время, общаясь с попутчиками`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[92] = &KettelItemAnswerImpl{
		ID:      92,
		TraitID: 20,
		// nolint
		MaleQuestion: `В ситуациях, которые могут стать опасными, я громко разговариваю, хотя это выглядит невежливо и нарушает спокойствие:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[93] = &KettelItemAnswerImpl{
		ID:      93,
		TraitID: 21,
		// nolint
		MaleQuestion: `Когда я чувствую холодное отношение со стороны знакомых людей, это обычно:`,
		Variants: []string{
			`не сильно влияет на моё состояние`,
			`зависит от ситуации`,
			`заметно задевает меня`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[94] = &KettelItemAnswerImpl{
		ID:      94,
		TraitID: 21,
		// nolint
		MaleQuestion: `Я смущаюсь, когда меня хвалят или говорят мне комплименты:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[95] = &KettelItemAnswerImpl{
		ID:      95,
		TraitID: 22,
		// nolint
		MaleQuestion: `Мне обычно комфортнее работа, где:`,
		Variants: []string{
			`доход достаточно стабилен и предсказуем`,
			`важны оба варианта`,
			`уровень дохода сильнее зависит от личных результатов`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[96] = &KettelItemAnswerImpl{
		ID:      96,
		TraitID: 23,
		// nolint
		MaleQuestion: `Когда мне нужно разобраться в чём-то новом, я скорее:`,
		Variants: []string{
			`спрошу у человека, который хорошо в этом разбирается`,
			`использую оба подхода`,
			`сам найду и изучу информацию`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[97] = &KettelItemAnswerImpl{
		ID:      97,
		TraitID: 23,
		// nolint
		MaleQuestion: `Мне нравится принимать активное участие в общественных проектах и волонтёрстве:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[98] = &KettelItemAnswerImpl{
		ID:      98,
		TraitID: 24,
		// nolint
		MaleQuestion: `Уже после завершения задачи я иногда продолжаю замечать мелкие вещи, которые можно было бы поправить:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[99] = &KettelItemAnswerImpl{
		ID:      99,
		TraitID: 25,
		// nolint
		MaleQuestion: `Мелкие накладки или ошибки иногда дольше остаются у меня в голове, чем хотелось бы:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[100] = &KettelItemAnswerImpl{
		ID:      100,
		TraitID: 25,
		// nolint
		MaleQuestion: `Обычно мне удаётся достаточно быстро расслабиться и спокойно уснуть:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[101] = &KettelItemAnswerImpl{
		ID:      101,
		TraitID: 10,
		// nolint
		MaleQuestion: `В работе мне обычно интереснее задачи, где:`,
		Variants: []string{
			`много общения и взаимодействия с людьми`,
			`важны оба варианта`,
			`нужно глубоко разбираться в информации и деталях`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[102] = &KettelItemAnswerImpl{
		ID:      102,
		TraitID: 11,
		// nolint
		MaleQuestion: `«Дирижёр» относится к «оркестру», как «режиссёр» относится к:`,
		Variants: []string{
			`«сценарию»`,
			`«театру»`,
			`«съёмочной группе»`,
		},
		RawVariantKeys: []int{0, 0, 1},
	}

	ItemsLib[103] = &KettelItemAnswerImpl{
		ID:      103,
		TraitID: 11,
		// nolint
		MaleQuestion: `«24» так относится к «53», как «68» относится к:`,
		Variants: []string{
			`«97»`,
			`«79»`,
			`«86»`,
		},
		RawVariantKeys: []int{0, 1, 0},
	}

	ItemsLib[104] = &KettelItemAnswerImpl{
		ID:      104,
		TraitID: 12,
		// nolint
		MaleQuestion: `Когда люди ведут себя неразумно, то я:`,
		Variants: []string{
			`молчу`,
			`не уверен`,
			`высказываю свое презрение`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[105] = &KettelItemAnswerImpl{
		ID:      105,
		TraitID: 12,
		// nolint
		MaleQuestion: `Если кто-нибудь громко разговаривает, когда я слушаю музыку:`,
		Variants: []string{
			`могу сосредоточиться на музыке, не отвлекаться`,
			`нечто среднее`,
			`чувствую, что это портит мне удовольствие и раздражает`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[106] = &KettelItemAnswerImpl{
		ID:      106,
		TraitID: 13,
		// nolint
		MaleQuestion: `Меня лучше охарактеризовать как:`,
		Variants: []string{
			`вежливого и спокойного`,
			`нечто среднее`,
			`энергичного и активного`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[107] = &KettelItemAnswerImpl{
		ID:      107,
		TraitID: 14,
		// nolint
		MaleQuestion: `В общественных мероприятиях я принимаю участие только тогда, когда это нужно, а в иных случаях избегаю их:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[108] = &KettelItemAnswerImpl{
		ID:      108,
		TraitID: 14,
		// nolint
		MaleQuestion: `Быть осторожным и не ждать хорошего лучше, чем быть оптимистом и всегда ждать успеха:`,
		// nolint
		FemaleQuestion: `Быть осторожной и не ждать хорошего лучше, чем быть оптимисткой и всегда ждать успеха:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[109] = &KettelItemAnswerImpl{
		ID:      109,
		TraitID: 15,
		// nolint
		MaleQuestion: `Думая о трудностях в своей работе, я:`,
		Variants: []string{
			`стараюсь планировать заранее, прежде чем встретить трудность`,
			`нечто среднее`,
			`считаю, что справлюсь с трудностями по мере того, как они возникнут`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[110] = &KettelItemAnswerImpl{
		ID:      110,
		TraitID: 16,
		// nolint
		MaleQuestion: `Мне легко вступить в контакт с людьми во время различных общественных мероприятий:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[111] = &KettelItemAnswerImpl{
		ID:      111,
		TraitID: 16,
		// nolint
		MaleQuestion: `Когда требуется немного дипломатии и умения убедить, чтобы побудить людей что-либо сделать, обычно об этом просят меня:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[112] = &KettelItemAnswerImpl{
		ID:      112,
		TraitID: 17,
		// nolint
		MaleQuestion: `Интересно быть:`,
		Variants: []string{
			`ментором, помогающим людям с карьерой`,
			`нечто среднее`,
			`руководителем IT-проекта`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[113] = &KettelItemAnswerImpl{
		ID:      113,
		TraitID: 18,
		// nolint
		MaleQuestion: `Если я уверен, что человек несправедлив или ведет себя эгоистично, я указываю на это, даже если это связано с неприятностями:`,
		// nolint
		FemaleQuestion: `Если я уверена, что человек несправедлив или ведет себя эгоистично, я указываю на это, даже если это связано с неприятностями:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[114] = &KettelItemAnswerImpl{
		ID:      114,
		TraitID: 18,
		// nolint
		MaleQuestion: `Иногда я специально говорю что-то странное или нелепое ради реакции окружающих:`,
		Variants: []string{
			`это про меня`,
			`бывает по-разному`,
			`это не про меня`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[115] = &KettelItemAnswerImpl{
		ID:      115,
		TraitID: 19,
		// nolint
		MaleQuestion: `Я мог бы получать удовольствие от работы, связанной с анализом театра, музыки и культурных событий:`,
		Variants: []string{
			`скорее да`,
			`затрудняюсь ответить`,
			`скорее нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[116] = &KettelItemAnswerImpl{
		ID:      116,
		TraitID: 19,
		// nolint
		MaleQuestion: `Когда мне приходится долго сидеть без движения, у меня обычно появляется желание чем-нибудь занять руки или отвлечься:`,
		Variants: []string{
			`часто так`,
			`по-разному`,
			`почти никогда`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[117] = &KettelItemAnswerImpl{
		ID:      117,
		TraitID: 20,
		// nolint
		MaleQuestion: `Если человек говорит что-то явно неверное, я скорее решу, что:`,
		Variants: []string{
			`он намеренно вводит в заблуждение`,
			`зависит от ситуации`,
			`ему просто не хватает информации`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[118] = &KettelItemAnswerImpl{
		ID:      118,
		TraitID: 21,
		// nolint
		MaleQuestion: `Я чувствую, что мне угрожает какое-то наказание, даже когда я ничего плохого не сделал:`,
		// nolint
		FemaleQuestion: `Я чувствую, что мне угрожает какое-то наказание, даже когда я ничего плохого не сделала:`,
		Variants: []string{
			`часто`,
			`иногда`,
			`никогда`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[119] = &KettelItemAnswerImpl{
		ID:      119,
		TraitID: 21,
		// nolint
		MaleQuestion: `Мнение о том, что болезнь также часто бывает от психических, как и от физических факторов, сильно преувеличено:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[120] = &KettelItemAnswerImpl{
		ID:      120,
		TraitID: 22,
		// nolint
		MaleQuestion: `Даже в современном мире традиционные ритуалы и торжественные церемонии не теряют своей ценности:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[121] = &KettelItemAnswerImpl{
		ID:      121,
		TraitID: 23,
		// nolint
		MaleQuestion: `Мысль о том, что люди подумают, будто я веду себя необычно или странно, меня беспокоит:`,
		Variants: []string{
			`очень`,
			`немного`,
			`совсем не беспокоит`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[122] = &KettelItemAnswerImpl{
		ID:      122,
		TraitID: 23,
		// nolint
		MaleQuestion: `Выполняя какое-либо дело, я бы предпочел работать:`,
		// nolint
		FemaleQuestion: `Выполняя какое-либо дело, я бы предпочла работать:`,
		Variants: []string{
			`в составе коллектива`,
			`не уверен`,
			`самостоятельно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[123] = &KettelItemAnswerImpl{
		ID:      123,
		TraitID: 24,
		// nolint
		MaleQuestion: `Иногда я надолго застреваю в переживаниях о том, как несправедливо со мной обошлись:`,
		Variants: []string{
			`часто`,
			`иногда`,
			`почти никогда`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[124] = &KettelItemAnswerImpl{
		ID:      124,
		TraitID: 25,
		// nolint
		MaleQuestion: `Иногда люди раздражают меня быстрее, чем стоило бы:`,
		Variants: []string{
			`часто`,
			`по-разному`,
			`редко`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[125] = &KettelItemAnswerImpl{
		ID:      125,
		TraitID: 25,
		// nolint
		MaleQuestion: `Я всегда могу без труда изменить свои старые привычки и не возвращаться к прежнему:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[126] = &KettelItemAnswerImpl{
		ID:      126,
		TraitID: 10,
		// nolint
		MaleQuestion: `При одинаковой оплате я скорее выбрал бы работу:`,
		Variants: []string{
			`где важны аргументы, правила и переговоры`,
			`что-то среднее`,
			`где есть техника, движение и ощущение свободы`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[127] = &KettelItemAnswerImpl{
		ID:      127,
		TraitID: 11,
		// nolint
		MaleQuestion: `«Лучшее» так относится к «наихудшее», как «медленное» к:`,
		Variants: []string{
			`«быстрое»`,
			`«лучшее»`,
			`«быстрейшее»`,
		},
		RawVariantKeys: []int{0, 0, 1},
	}

	ItemsLib[128] = &KettelItemAnswerImpl{
		ID:      128,
		TraitID: 11,
		// nolint
		MaleQuestion: `Каким из приведенных ниже сочетаний следует продолжить буквенный ряд РООООРРОООРРР...:`,
		Variants: []string{
			`ОРРР`,
			`ООРР`,
			`РООО`,
		},
		RawVariantKeys: []int{0, 1, 0},
	}

	ItemsLib[129] = &KettelItemAnswerImpl{
		ID:      129,
		TraitID: 12,
		// nolint
		MaleQuestion: `Когда цель становится реально достижимой, мой интерес к ней иногда резко снижается:`,
		// nolint
		FemaleQuestion: `Когда цель становится реально достижимой, мой интерес к ней иногда резко снижается:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[130] = &KettelItemAnswerImpl{
		ID:      130,
		TraitID: 12,
		// nolint
		MaleQuestion: `Большей частью я могу продолжать работать тщательно, не обращая внимания на шум, создаваемый другими:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[131] = &KettelItemAnswerImpl{
		ID:      131,
		TraitID: 13,
		// nolint
		MaleQuestion: `Иногда я говорю посторонним вещи, кажущиеся мне важными, независимо от того, спрашивают ли они об этом:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[132] = &KettelItemAnswerImpl{
		ID:      132,
		TraitID: 14,
		// nolint
		MaleQuestion: `Мы с друзьями часто вспоминаем и обсуждаем истории, которые пережили вместе:`,
		// nolint
		FemaleQuestion: `Мы с друзьями часто вспоминаем и обсуждаем истории, которые пережили вместе:`,
		Variants: []string{
			`да`,
			`иногда`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[133] = &KettelItemAnswerImpl{
		ID:      133,
		TraitID: 14,
		// nolint
		MaleQuestion: `Иногда мне нравится шутить на грани или делать что-то ради острых ощущений и реакции других:`,
		Variants: []string{
			`это про меня`,
			`иногда`,
			`это не про меня`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[134] = &KettelItemAnswerImpl{
		ID:      134,
		TraitID: 15,
		// nolint
		MaleQuestion: `Вид неубранной комнаты очень раздражает меня:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[135] = &KettelItemAnswerImpl{
		ID:      135,
		TraitID: 16,
		// nolint
		MaleQuestion: `Я считаю себя общительным открытым человеком:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[136] = &KettelItemAnswerImpl{
		ID:      136,
		TraitID: 16,
		// nolint
		MaleQuestion: `В общении я:`,
		Variants: []string{
			`свободно проявляю свои чувства`,
			`нечто среднее`,
			`держу свои переживания «при себе»`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[137] = &KettelItemAnswerImpl{
		ID:      137,
		TraitID: 17,
		// nolint
		MaleQuestion: `Я люблю музыку:`,
		Variants: []string{
			`скорее энергичную и ритмичную`,
			`нечто среднее`,
			`скорее глубокую и эмоциональную`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[138] = &KettelItemAnswerImpl{
		ID:      138,
		TraitID: 17,
		// nolint
		MaleQuestion: `Я скорее оценю продукт за то, как он решает задачи пользователей, чем за элегантность кода и инженерных решений внутри:`,
		Variants: []string{
			`это про меня`,
			`частично`,
			`это не про меня`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[139] = &KettelItemAnswerImpl{
		ID:      139,
		TraitID: 18,
		// nolint
		MaleQuestion: `Если моя удачная шутка остается незамеченной окружающими, то я:`,
		Variants: []string{
			`оставлю это`,
			`нечто среднее`,
			`дам людям возможность услышать её ещё раз`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[140] = &KettelItemAnswerImpl{
		ID:      140,
		TraitID: 19,
		// nolint
		MaleQuestion: `Мне была бы интересна работа, связанная с созданием визуального контента и съёмкой для брендов или соцсетей:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[141] = &KettelItemAnswerImpl{
		ID:      141,
		TraitID: 19,
		// nolint
		MaleQuestion: `С незнакомыми людьми лучше держать дистанцию — никогда не знаешь, чего от них ожидать:`,
		// nolint
		FemaleQuestion: `С незнакомыми людьми лучше держать дистанцию — никогда не знаешь, чего от них ожидать:`,
		Variants: []string{
			`скорее согласен`,
			`затрудняюсь ответить`,
			`скорее не согласен`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[142] = &KettelItemAnswerImpl{
		ID:      142,
		TraitID: 20,
		// nolint
		MaleQuestion:   `Путешествуя в новой стране, я скорее выберу готовый маршрут или сопровождение, чем буду всё организовывать самостоятельно:`,
		FemaleQuestion: `Путешествуя в новой стране, я скорее выберу готовый маршрут или сопровождение, чем буду всё организовывать самостоятельно:`,
		Variants: []string{
			`скорее да`,
			`не уверен`,
			`скорее нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[143] = &KettelItemAnswerImpl{
		ID:      143,
		TraitID: 21,
		// nolint
		MaleQuestion: `Я часто чувствую, что прикладываю много усилий, но результаты оказываются скромнее, чем хотелось бы:`,
		// nolint
		FemaleQuestion: `Я часто чувствую, что прикладываю много усилий, но результаты оказываются скромнее, чем хотелось бы:`,
		Variants: []string{
			`скорее да`,
			`иногда`,
			`скорее нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[144] = &KettelItemAnswerImpl{
		ID:      144,
		TraitID: 21,
		// nolint
		MaleQuestion: `Если кто-то пользуется моей добротой или уступчивостью, я обычно не держу на него долгой обиды:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[145] = &KettelItemAnswerImpl{
		ID:      145,
		TraitID: 22,
		// nolint
		MaleQuestion: `Когда спор становится слишком эмоциональным и жёстким, я скорее хочу:`,
		Variants: []string{
			`понять, чья позиция окажется сильнее`,
			`зависит от ситуации`,
			`чтобы люди спокойно договорились`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[146] = &KettelItemAnswerImpl{
		ID:      146,
		TraitID: 23,
		// nolint
		MaleQuestion: `Я предпочитаю планировать что-либо самостоятельно, без вмешательства и предложений со стороны:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[147] = &KettelItemAnswerImpl{
		ID:      147,
		TraitID: 24,
		// nolint
		MaleQuestion: `Иногда чувство зависти влияет на мои действия:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[148] = &KettelItemAnswerImpl{
		ID:      148,
		TraitID: 24,
		// nolint
		MaleQuestion: `Я считаю, что авторитет руководителя не должен полностью зависеть от того, согласны ли с ним окружающие:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[149] = &KettelItemAnswerImpl{
		ID:      149,
		TraitID: 25,
		// nolint
		MaleQuestion: `Список предстоящих дел обычно держит меня в тонусе:`,
		Variants: []string{
			`да`,
			`иногда`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[150] = &KettelItemAnswerImpl{
		ID:      150,
		TraitID: 25,
		// nolint
		MaleQuestion: `Когда окружающие эмоционально комментируют мои действия или результаты, я обычно сохраняю спокойствие:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[151] = &KettelItemAnswerImpl{
		ID:      151,
		TraitID: 10,
		// nolint
		MaleQuestion: `Я скорее получил бы удовольствие от работы, где нужно:`,
		Variants: []string{
			`продумывать интерфейсы и опыт пользователей`,
			`что-то среднее`,
			`координировать события и взаимодействие людей`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[152] = &KettelItemAnswerImpl{
		ID:      152,
		TraitID: 11,
		// nolint
		MaleQuestion: `Какое слово выбивается из группы:`,
		Variants: []string{
			`импровизация`,
			`структура`,
			`план`,
		},
		RawVariantKeys: []int{1, 0, 0},
	}

	ItemsLib[153] = &KettelItemAnswerImpl{
		ID:      153,
		TraitID: 11,
		// nolint
		MaleQuestion: `«Пламя» так относится к «жар», как «роза» относится к:`,
		Variants: []string{
			`«шип»`,
			`«красивые лепестки»`,
			`«аромат»`,
		},
		RawVariantKeys: []int{0, 0, 1},
	}

	ItemsLib[154] = &KettelItemAnswerImpl{
		ID:      154,
		TraitID: 12,
		// nolint
		MaleQuestion: `У меня бывают яркие сновидения, мешающие мне спать:`,
		Variants: []string{
			`часто`,
			`иногда`,
			`практически никогда`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[155] = &KettelItemAnswerImpl{
		ID:      155,
		TraitID: 13,
		// nolint
		MaleQuestion: `Если на пути к успеху стоят серьезные препятствия, я все-таки предпочитаю рискнуть:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[156] = &KettelItemAnswerImpl{
		ID:      156,
		TraitID: 13,
		// nolint
		MaleQuestion: `Когда я нахожусь в группе людей, приступающих к какой-то работе, то само собой получается, что я оказываюсь во главе их:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[157] = &KettelItemAnswerImpl{
		ID:      157,
		TraitID: 14,
		// nolint
		MaleQuestion: `Мне больше нравится в одежде спокойный стиль, чем бросающаяся в глаза индивидуальность:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[158] = &KettelItemAnswerImpl{
		ID:      158,
		TraitID: 14,
		// nolint
		MaleQuestion: `Мне больше нравится провести вечер за спокойным хобби, чем в шумной компании:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[159] = &KettelItemAnswerImpl{
		ID:      159,
		TraitID: 15,
		// nolint
		MaleQuestion: `Бывает, что я не придаю значения советам других людей, даже если позже оказывается, что они были полезными:`,
		Variants: []string{
			`иногда`,
			`почти никогда`,
			`никогда`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[160] = &KettelItemAnswerImpl{
		ID:      160,
		TraitID: 15,
		// nolint
		MaleQuestion: `В своих поступках я всегда стараюсь придерживаться общепринятых правил поведения:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[161] = &KettelItemAnswerImpl{
		ID:      161,
		TraitID: 16,
		// nolint
		MaleQuestion: `Мне не очень нравится, когда смотрят, как я работаю:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[162] = &KettelItemAnswerImpl{
		ID:      162,
		TraitID: 17,
		// nolint
		MaleQuestion: `Иногда приходится применять силу, потому что не всегда возможно добиться результата с помощью утверждения:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[163] = &KettelItemAnswerImpl{
		ID:      163,
		TraitID: 17,
		// nolint
		MaleQuestion: `Во время учёбы мне больше нравились:`,
		// nolint
		FemaleQuestion: `Во время учёбы мне больше нравились:`,
		Variants: []string{
			`сочинения, тексты и гуманитарные предметы`,
			`нечто среднее`,
			`математика и логические задачи`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[164] = &KettelItemAnswerImpl{
		ID:      164,
		TraitID: 18,
		// nolint
		MaleQuestion: `Меня иногда огорчало, что обо мне за глаза отзывались неодобрительно без всяких к этому причин:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[165] = &KettelItemAnswerImpl{
		ID:      165,
		TraitID: 19,
		// nolint
		MaleQuestion: `Разговор с простыми людьми, которые всегда придерживаются общепринятых правил и традиций:`,
		Variants: []string{
			`часто вполне интересен и содержателен`,
			`нечто среднее`,
			`раздражает меня, потому что ограничивается мелочами`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[166] = &KettelItemAnswerImpl{
		ID:      166,
		TraitID: 19,
		// nolint
		MaleQuestion: `Некоторые вещи настолько раздражают меня, что предпочитаю вообще не говорить на эти темы:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[167] = &KettelItemAnswerImpl{
		ID:      167,
		TraitID: 20,
		// nolint
		MaleQuestion: `В воспитании важнее:`,
		Variants: []string{
			`относиться к ребёнку с достаточной любовью`,
			`нечто среднее`,
			`сформировать правильные привычки и отношение к жизни`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[168] = &KettelItemAnswerImpl{
		ID:      168,
		TraitID: 21,
		// nolint
		MaleQuestion: `Окружающие обычно воспринимают меня как спокойного и уравновешенного человека, которого трудно выбить из колеи:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[169] = &KettelItemAnswerImpl{
		ID:      169,
		TraitID: 22,
		// nolint
		MaleQuestion: `Мне ближе идея, что правила и общественные нормы стоит пересматривать, если они уже не соответствуют современности:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[170] = &KettelItemAnswerImpl{
		ID:      170,
		TraitID: 22,
		// nolint
		MaleQuestion: `Думаю, что в современном мире важнее разрешить:`,
		Variants: []string{
			`вопросы нравственности и этики`,
			`не уверен`,
			`глобальные конфликты между странами`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[171] = &KettelItemAnswerImpl{
		ID:      171,
		TraitID: 23,
		// nolint
		MaleQuestion: `Мне проще разобраться в новой теме:`,
		Variants: []string{
			`самостоятельно изучая хороший материал`,
			`зависит от ситуации`,
			`обсуждая её с другими людьми`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[172] = &KettelItemAnswerImpl{
		ID:      172,
		TraitID: 24,
		// nolint
		MaleQuestion: `Я предпочитаю идти своим путем вместо того, чтобы действовать в соответствии с принятыми правилами:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[173] = &KettelItemAnswerImpl{
		ID:      173,
		TraitID: 24,
		// nolint
		MaleQuestion: `Прежде чем выдвигать какой-либо аргумент, я предпочитаю подождать, пока не буду убежден, что я прав:`,
		// nolint
		FemaleQuestion: `Прежде чем выдвигать какой-либо аргумент, я предпочитаю подождать, пока не буду убеждена, что я права:`,
		Variants: []string{
			`всегда`,
			`обычно`,
			`только если это целесообразно`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[174] = &KettelItemAnswerImpl{
		ID:      174,
		TraitID: 25,
		// nolint
		MaleQuestion: `Мелочи иногда невыносимо «действуют мне на нервы», хотя я и понимаю, что они не существенны:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[175] = &KettelItemAnswerImpl{
		ID:      175,
		TraitID: 25,
		// nolint
		MaleQuestion: `У меня редко бывает, что в порыве эмоций я говорю то, о чём потом жалею:`,
		Variants: []string{
			`это про меня`,
			`затрудняюсь ответить`,
			`это не про меня`,
		},
		RawVariantKeys: []int{0, 1, 2},
	}

	ItemsLib[176] = &KettelItemAnswerImpl{
		ID:      176,
		TraitID: 10,
		// nolint
		MaleQuestion: `Если бы меня попросили участвовать в волонтёрской деятельности, то я бы:`,
		Variants: []string{
			`согласился`,
			`не уверен`,
			`вежливо сказал, что занят`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[177] = &KettelItemAnswerImpl{
		ID:      177,
		TraitID: 11,
		// nolint
		MaleQuestion: `Которое из следующих слов не относится к двум другим:`,
		Variants: []string{
			`широкий`,
			`зигзагообразный`,
			`прямой`,
		},
		RawVariantKeys: []int{1, 0, 0},
	}

	ItemsLib[178] = &KettelItemAnswerImpl{
		ID:      178,
		TraitID: 11,
		// nolint
		MaleQuestion: `«Скоро» так относится к «никогда», как «близко» к:`,
		Variants: []string{
			`«нигде»`,
			`«далеко»`,
			`«где-то»`,
		},
		RawVariantKeys: []int{1, 0, 0},
	}

	ItemsLib[179] = &KettelItemAnswerImpl{
		ID:      179,
		TraitID: 12,
		// nolint
		MaleQuestion: `Если я невольно нарушил правила поведения, находясь в обществе, то я вскоре забываю об этом:`,
		// nolint
		FemaleQuestion: `Если я невольно нарушила правила поведения, находясь в обществе, то я вскоре забываю об этом:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[180] = &KettelItemAnswerImpl{
		ID:      180,
		TraitID: 13,
		// nolint
		MaleQuestion: `Когда возникает проблема, люди нередко ждут от меня полезных идей или нестандартного решения:`,
		Variants: []string{
			`скорее да`,
			`по-разному`,
			`скорее нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[181] = &KettelItemAnswerImpl{
		ID:      181,
		TraitID: 13,
		// nolint
		MaleQuestion: `Я способен лучше проявить себя:`,
		// nolint
		FemaleQuestion: `Я способна лучше проявить себя:`,
		Variants: []string{
			`в трудных ситуациях, когда нужно сохранить самообладание`,
			`не уверен`,
			`когда требуется умение ладить с людьми`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[182] = &KettelItemAnswerImpl{
		ID:      182,
		TraitID: 14,
		// nolint
		MaleQuestion: `Меня считают человеком, полным энтузиазма и энергии:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[183] = &KettelItemAnswerImpl{
		ID:      183,
		TraitID: 14,
		// nolint
		MaleQuestion: `Мне нравится работа, которая требует перемен, разнообразия, командировок, даже если связана с некоторой опасностью:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[184] = &KettelItemAnswerImpl{
		ID:      184,
		TraitID: 15,
		// nolint
		MaleQuestion: `Я довольно требовательный человек и всегда настаиваю на том, чтобы все делалось по возможности правильно:`,
		// nolint
		FemaleQuestion: `Я довольно требовательная и всегда настаиваю на том, чтобы все делалось по возможности правильно:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[185] = &KettelItemAnswerImpl{
		ID:      185,
		TraitID: 15,
		// nolint
		MaleQuestion: `Мне нравится работа, требующая внимательного отношения, точных навыков:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}

	ItemsLib[186] = &KettelItemAnswerImpl{
		ID:      186,
		TraitID: 16,
		// nolint
		MaleQuestion: `Я отношусь к активным людям, которые всегда чем-то заняты:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
		RawVariantKeys: []int{2, 1, 0},
	}
}
