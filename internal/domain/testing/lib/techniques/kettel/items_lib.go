package kettel

var ItemsLib = map[uint64]*KettelItemAnswerImpl{}

func init() {
	ItemsLib[3] = &KettelItemAnswerImpl{
		ID:      3,
		TraitID: 10,
		// nolint
		MaleQuestion: `Я бы предпочел временами жить в доме, который находится:`,
		// nolint
		FemaleQuestion: `Я бы предпочла временами жить в доме, который находится:`,
		Variants: []string{
			`в обжитом городе`,
			`нечто среднее`,
			`одиноко в глухих лесах`,
		},
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
	}

	ItemsLib[8] = &KettelItemAnswerImpl{
		ID:      8,
		TraitID: 14,
		// nolint
		MaleQuestion: `Мне больше нравится классическая, чем эстрадная музыка:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
	}

	ItemsLib[9] = &KettelItemAnswerImpl{
		ID:      9,
		TraitID: 15,
		// nolint
		MaleQuestion: `Если бы я увидел дерущимися соседских детей, то я:`,
		// nolint
		FemaleQuestion: `Если бы я увидела дерущимися соседских детей, то я:`,
		Variants: []string{
			`дал бы им возможность договориться самим`,
			`не уверен`,
			`рассудил бы их`,
		},
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
	}

	ItemsLib[11] = &KettelItemAnswerImpl{
		ID:      11,
		TraitID: 17,
		// nolint
		MaleQuestion: `По-моему, интереснее быть:`,
		Variants: []string{
			`инженером`,
			`не уверен`,
			`журналистом`,
		},
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
	}

	ItemsLib[15] = &KettelItemAnswerImpl{
		ID:      15,
		TraitID: 19,
		// nolint
		MaleQuestion: `Было бы хорошо, если бы отпуск (каникулы) был более продолжителен, и каждый был бы обязан его использовать:`,
		Variants: []string{
			`согласен`,
			`не уверен`,
			`не согласен`,
		},
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
	}

	ItemsLib[17] = &KettelItemAnswerImpl{
		ID:      17,
		TraitID: 20,
		// nolint
		MaleQuestion: `Я говорю о своих чувствах:`,
		Variants: []string{
			`только если это необходимо`,
			`нечто среднее`,
			`охотно, когда представится возможность`,
		},
	}

	ItemsLib[18] = &KettelItemAnswerImpl{
		ID:      18,
		TraitID: 21,
		// nolint
		MaleQuestion: `Время от времени у меня возникает чувство неопределенной опасности или внезапного страха по непонятным причинам:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
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
	}

	ItemsLib[21] = &KettelItemAnswerImpl{
		ID:      21,
		TraitID: 22,
		// nolint
		MaleQuestion: `Моим решением руководит больше:`,
		Variants: []string{
			`сердце`,
			`сердце и разум в равной степени`,
			`разум`,
		},
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
	}

	ItemsLib[23] = &KettelItemAnswerImpl{
		ID:      23,
		TraitID: 24,
		// nolint
		MaleQuestion: `Иногда, когда я смотрю в зеркало, мне трудно разобраться, где у меня правая, а где левая сторона:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
	}

	ItemsLib[24] = &KettelItemAnswerImpl{
		ID:      24,
		TraitID: 24,
		// nolint
		MaleQuestion: `При разговоре я предпочитаю:`,
		Variants: []string{
			`высказывать свои мысли так, как они приходят мне в голову`,
			`нечто среднее`,
			`сначала сформулировать получше свои мысли`,
		},
	}

	ItemsLib[25] = &KettelItemAnswerImpl{
		ID:      25,
		TraitID: 25,
		// nolint
		MaleQuestion: `После того как меня что-то сильно рассердит, я довольно быстро успокаиваюсь:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[26] = &KettelItemAnswerImpl{
		ID:      26,
		TraitID: 10,
		// nolint
		MaleQuestion: `При одинаковом рабочем времени и заработке было бы интереснее работать:`,
		Variants: []string{
			`плотником или поваром`,
			`не уверен`,
			`официантом в хорошем ресторане`,
		},
	}

	ItemsLib[27] = &KettelItemAnswerImpl{
		ID:      27,
		TraitID: 10,
		// nolint
		MaleQuestion: `На общественные должности меня выбирали:`,
		Variants: []string{
			`очень редко`,
			`иногда`,
			`много раз`,
		},
	}

	ItemsLib[28] = &KettelItemAnswerImpl{
		ID:      28,
		TraitID: 11,
		// nolint
		MaleQuestion: `«Лопата» относится к «копать», как «нож» относится к:`,
		Variants: []string{
			`«острый»`,
			`«резать»`,
			`«указывать»`,
		},
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
	}

	ItemsLib[30] = &KettelItemAnswerImpl{
		ID:      30,
		TraitID: 12,
		// nolint
		MaleQuestion: `В своей жизни я почти всегда достигаю поставленных целей:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
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
	}

	ItemsLib[32] = &KettelItemAnswerImpl{
		ID:      32,
		TraitID: 13,
		// nolint
		MaleQuestion: `Я чувствую себя «не в своей тарелке», когда мне приходится работать над чем-нибудь, что требует быстрых действий, результаты которых могут повлиять на других людей:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
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
	}

	ItemsLib[34] = &KettelItemAnswerImpl{
		ID:      34,
		TraitID: 15,
		// nolint
		MaleQuestion: `Когда я вижу неряшливых, неопрятных людей, я:`,
		Variants: []string{
			`принимаю их такими, как они есть`,
			`нечто среднее`,
			`испытываю отвращение и возмущение`,
		},
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
	}

	ItemsLib[36] = &KettelItemAnswerImpl{
		ID:      36,
		TraitID: 16,
		// nolint
		MaleQuestion: `Я всегда рад оказаться среди людей, например, в гостях, на танцах, коллективной встрече:`,
		// nolint
		FemaleQuestion: `Я всегда рада оказаться среди людей, например, в гостях, на танцах, коллективной встрече:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[37] = &KettelItemAnswerImpl{
		ID:      37,
		TraitID: 17,
		// nolint
		MaleQuestion: `В школе я предпочитал (или предпочитаю):`,
		// nolint
		FemaleQuestion: `В школе я предпочитала (или предпочитаю):`,
		Variants: []string{
			`заниматься музыкой, пением`,
			`нечто среднее`,
			`выпиливать и мастерить что-либо`,
		},
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
	}

	ItemsLib[41] = &KettelItemAnswerImpl{
		ID:      41,
		TraitID: 20,
		// nolint
		MaleQuestion: `Время от времени у меня появляется потребность в интересной физической деятельности:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[42] = &KettelItemAnswerImpl{
		ID:      42,
		TraitID: 20,
		// nolint
		MaleQuestion: `Я предпочел бы скорее общаться с вежливыми людьми, чем с грубоватыми и любящими возражать:`,
		// nolint
		FemaleQuestion: `Я предпочла бы скорее общаться с вежливыми людьми, чем с грубоватыми и любящими возражать:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
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
	}

	ItemsLib[44] = &KettelItemAnswerImpl{
		ID:      44,
		TraitID: 21,
		// nolint
		MaleQuestion: `Если меня вызывает начальство, то я:`,
		Variants: []string{
			`пользуюсь случаем, чтобы попросить о чем-то нужном мне`,
			`нечто среднее`,
			`боюсь, что это связано с какой-нибудь оплошностью в моей работе`,
		},
	}

	ItemsLib[45] = &KettelItemAnswerImpl{
		ID:      45,
		TraitID: 22,
		// nolint
		MaleQuestion: `В наше время требуется:`,
		Variants: []string{
			`больше спокойных, солидных людей`,
			`не уверен`,
			`больше «идеалистов», планирующих лучшее будущее`,
		},
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
	}

	ItemsLib[47] = &KettelItemAnswerImpl{
		ID:      47,
		TraitID: 23,
		// nolint
		MaleQuestion: `В юности я принимал участие в нескольких спортивных мероприятиях:`,
		// nolint
		FemaleQuestion: `В юности я принимала участие в нескольких спортивных мероприятиях:`,
		Variants: []string{
			`иногда`,
			`довольно часто`,
			`многократно`,
		},
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
	}

	ItemsLib[51] = &KettelItemAnswerImpl{
		ID:      51,
		TraitID: 10,
		// nolint
		MaleQuestion: `Если бы пришлось выбирать, то я предпочел бы быть:`,
		// nolint
		FemaleQuestion: `Если бы пришлось выбирать, то я предпочла бы быть:`,
		Variants: []string{
			`лесником`,
			`не уверен`,
			`учителем средней школы`,
		},
	}

	ItemsLib[52] = &KettelItemAnswerImpl{
		ID:      52,
		TraitID: 10,
		// nolint
		MaleQuestion: `На праздники и дни рождения я:`,
		Variants: []string{
			`люблю делать подарки`,
			`неопределенно`,
			`считаю, что делать подарки – довольно неприятная вещь`,
		},
	}

	ItemsLib[53] = &KettelItemAnswerImpl{
		ID:      53,
		TraitID: 11,
		// nolint
		MaleQuestion: `«Усталый» относится к «работе», как «гордый» к:`,
		Variants: []string{
			`«улыбка»`,
			`«успех»`,
			`«счастливый»`,
		},
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
	}

	ItemsLib[58] = &KettelItemAnswerImpl{
		ID:      58,
		TraitID: 14,
		// nolint
		MaleQuestion: `Я склонен посещать зрелищные мероприятия и развлечения:`,
		// nolint
		FemaleQuestion: `Я склонна посещать зрелищные мероприятия и развлечения:`,
		Variants: []string{
			`чаще, чем раз в неделю (т.е. чаще, чем большинство)`,
			`примерно раз в неделю (т.е. как большинство)`,
			`реже, чем раз в неделю (т.е. реже, чем большинство)`,
		},
	}

	ItemsLib[59] = &KettelItemAnswerImpl{
		ID:      59,
		TraitID: 15,
		// nolint
		MaleQuestion: `Я считаю, что возможность вести себя непринужденно важнее, чем хорошие манеры и уважение к существующим правилам поведения:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
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
	}

	ItemsLib[61] = &KettelItemAnswerImpl{
		ID:      61,
		TraitID: 16,
		// nolint
		MaleQuestion: `Мне трудно говорить или декламировать перед большой группой людей:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[62] = &KettelItemAnswerImpl{
		ID:      62,
		TraitID: 17,
		// nolint
		MaleQuestion: `У меня хорошее чувство ориентировки в незнакомом месте (мне легко сказать, где север – восток – юг – запад):`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
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
	}

	ItemsLib[64] = &KettelItemAnswerImpl{
		ID:      64,
		TraitID: 18,
		// nolint
		MaleQuestion: `Встречаясь с несправедливостью, я скорее склонен забыть об этом, чем реагировать:`,
		// nolint
		FemaleQuestion: `Встречаясь с несправедливостью, я скорее склонна забыть об этом, чем реагировать:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
	}

	ItemsLib[65] = &KettelItemAnswerImpl{
		ID:      65,
		TraitID: 19,
		// nolint
		MaleQuestion: `Из моей памяти часто выпадают несущественные тривиальные вещи, например, названия улиц, магазинов:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[66] = &KettelItemAnswerImpl{
		ID:      66,
		TraitID: 20,
		// nolint
		MaleQuestion: `Мне бы понравилась жизнь ветеринара, лечение и операции на животных:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[67] = &KettelItemAnswerImpl{
		ID:      67,
		TraitID: 20,
		// nolint
		MaleQuestion: `Я ем со вкусом, не всегда так аккуратно и тщательно как другие люди:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`неверно`,
		},
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
	}

	ItemsLib[69] = &KettelItemAnswerImpl{
		ID:      69,
		TraitID: 21,
		// nolint
		MaleQuestion: `Иногда меня предупреждают о том, что в моем голосе и манерах слишком проявляется возбуждение:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[70] = &KettelItemAnswerImpl{
		ID:      70,
		TraitID: 22,
		// nolint
		MaleQuestion: `В юности, если я расходился во мнении с родителями, то я:`,
		// nolint
		FemaleQuestion: `В юности, если я расходилась во мнении с родителями, то я:`,
		Variants: []string{
			`оставался при своем мнении`,
			`нечто среднее`,
			`соглашался с их авторитетом`,
		},
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
	}

	ItemsLib[73] = &KettelItemAnswerImpl{
		ID:      73,
		TraitID: 24,
		// nolint
		MaleQuestion: `В большинстве случаев я чувствую себя зрелым человеком:`,
		// nolint
		FemaleQuestion: `В большинстве случаев я чувствую себя зрелой женщиной:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
	}

	ItemsLib[74] = &KettelItemAnswerImpl{
		ID:      74,
		TraitID: 25,
		// nolint
		MaleQuestion: `Замечания в мой адрес, которые позволяют себе некоторые люди, меня больше расстраивают, чем помогают:`,
		Variants: []string{
			`часто`,
			`иногда`,
			`никогда`,
		},
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
	}

	ItemsLib[76] = &KettelItemAnswerImpl{
		ID:      76,
		TraitID: 10,
		// nolint
		MaleQuestion: `Начиная работу над полезным изобретением, я бы предпочел:`,
		// nolint
		FemaleQuestion: `Начиная работу над полезным изобретением, я бы предпочла:`,
		Variants: []string{
			`разрабатывать его в лаборатории`,
			`нечто среднее`,
			`заниматься его практической реализацией`,
		},
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
	}

	ItemsLib[79] = &KettelItemAnswerImpl{
		ID:      79,
		TraitID: 12,
		// nolint
		MaleQuestion: `Кажется, некоторые люди игнорируют и избегают меня, хотя я не знаю, почему:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
	}

	ItemsLib[80] = &KettelItemAnswerImpl{
		ID:      80,
		TraitID: 12,
		// nolint
		MaleQuestion: `Отношения ко мне людей не соответствуют моим добрым намерениям:`,
		Variants: []string{
			`часто`,
			`иногда`,
			`никогда`,
		},
	}

	ItemsLib[81] = &KettelItemAnswerImpl{
		ID:      81,
		TraitID: 13,
		// nolint
		MaleQuestion: `Употребление нецензурных выражений вызывает у меня возмущение, даже если не присутствуют лица другого пола:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
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
	}

	ItemsLib[85] = &KettelItemAnswerImpl{
		ID:      85,
		TraitID: 16,
		// nolint
		MaleQuestion: `Волнение перед выступлением в присутствии многих людей я испытывал:`,
		// nolint
		FemaleQuestion: `Волнение перед выступлением в присутствии многих людей я испытывала:`,
		Variants: []string{
			`довольно часто`,
			`иногда`,
			`почти никогда`,
		},
	}

	ItemsLib[86] = &KettelItemAnswerImpl{
		ID:      86,
		TraitID: 16,
		// nolint
		MaleQuestion: `Когда я нахожусь в большой группе людей, то я предпочитаю молчать и предоставляю слово другим:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[87] = &KettelItemAnswerImpl{
		ID:      87,
		TraitID: 17,
		// nolint
		MaleQuestion: `Я предпочитаю читать:`,
		Variants: []string{
			`реалистические описания военных и политических сражений`,
			`нечто среднее`,
			`роман, где много чувств и воображения`,
		},
	}

	ItemsLib[88] = &KettelItemAnswerImpl{
		ID:      88,
		TraitID: 18,
		// nolint
		MaleQuestion: `Когда люди пытаются мною командовать, то я поступаю как раз наоборот:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[89] = &KettelItemAnswerImpl{
		ID:      89,
		TraitID: 18,
		// nolint
		MaleQuestion: `Начальник или члены моей семьи критикуют меня только тогда, когда к этому действительно есть повод:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
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
	}

	ItemsLib[91] = &KettelItemAnswerImpl{
		ID:      91,
		TraitID: 19,
		// nolint
		MaleQuestion: `Во время длительной поездки я бы предпочел:`,
		// nolint
		FemaleQuestion: `Во время длительной поездки я бы предпочла:`,
		Variants: []string{
			`читать что-нибудь серьезное, но интересное`,
			`неопределенно`,
			`провести время, беседуя с кем-нибудь из пассажиров`,
		},
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
	}

	ItemsLib[93] = &KettelItemAnswerImpl{
		ID:      93,
		TraitID: 21,
		// nolint
		MaleQuestion: `Если знакомые плохо обращаются со мной и показывают свою неприязнь, то:`,
		Variants: []string{
			`меня это совершенно не трогает`,
			`нечто среднее`,
			`я расстраиваюсь`,
		},
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
	}

	ItemsLib[95] = &KettelItemAnswerImpl{
		ID:      95,
		TraitID: 22,
		// nolint
		MaleQuestion: `Я бы предпочел иметь работу:`,
		// nolint
		FemaleQuestion: `Я бы предпочла иметь работу:`,
		Variants: []string{
			`с постоянным окладом`,
			`нечто среднее`,
			`с большим окладом, который бы зависел от моей способности показать людям, чего я стою`,
		},
	}

	ItemsLib[96] = &KettelItemAnswerImpl{
		ID:      96,
		TraitID: 23,
		// nolint
		MaleQuestion: `Чтобы быть информированным, я предпочитаю получать сведения:`,
		// nolint
		FemaleQuestion: `Чтобы быть информированной, я предпочитаю получать сведения:`,
		Variants: []string{
			`в общении с людьми`,
			`нечто среднее`,
			`из литературы`,
		},
	}

	ItemsLib[97] = &KettelItemAnswerImpl{
		ID:      97,
		TraitID: 23,
		// nolint
		MaleQuestion: `Мне нравится принимать активное участие в общественной работе:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[98] = &KettelItemAnswerImpl{
		ID:      98,
		TraitID: 24,
		// nolint
		MaleQuestion: `При выполнении задания я удовлетворяюсь только тогда, когда должное внимание будет уделено всем мелочам:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
	}

	ItemsLib[99] = &KettelItemAnswerImpl{
		ID:      99,
		TraitID: 25,
		// nolint
		MaleQuestion: `Даже самые незначительные неудачи иногда меня слишком раздражают:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[100] = &KettelItemAnswerImpl{
		ID:      100,
		TraitID: 25,
		// nolint
		MaleQuestion: `Сон у меня всегда крепкий, я никогда не хожу и не разговариваю во сне:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[101] = &KettelItemAnswerImpl{
		ID:      101,
		TraitID: 10,
		// nolint
		MaleQuestion: `Для меня интереснее работа, при которой:`,
		Variants: []string{
			`нужно разговаривать с людьми`,
			`нечто среднее`,
			`нужно заниматься счетами и записями`,
		},
	}

	ItemsLib[102] = &KettelItemAnswerImpl{
		ID:      102,
		TraitID: 11,
		// nolint
		MaleQuestion: `«Размер» так относится к «длине», как «нечестный» к:`,
		Variants: []string{
			`«тюрьма»`,
			`«нарушение»`,
			`«кража»`,
		},
	}

	ItemsLib[103] = &KettelItemAnswerImpl{
		ID:      103,
		TraitID: 11,
		// nolint
		MaleQuestion: `«АБ» так относится к «ГВ», как «СР» относится к:`,
		Variants: []string{
			`«ПО»`,
			`«ОП»`,
			`«ТУ»`,
		},
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
	}

	ItemsLib[106] = &KettelItemAnswerImpl{
		ID:      106,
		TraitID: 13,
		// nolint
		MaleQuestion: `Меня лучше характеризовать как:`,
		Variants: []string{
			`вежливого и спокойного`,
			`нечто среднее`,
			`энергичного`,
		},
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
	}

	ItemsLib[112] = &KettelItemAnswerImpl{
		ID:      112,
		TraitID: 17,
		// nolint
		MaleQuestion: `Интересно быть:`,
		Variants: []string{
			`консультантом, помогающим людям выбирать профессию`,
			`нечто среднее`,
			`руководителем технического предприятия`,
		},
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
	}

	ItemsLib[114] = &KettelItemAnswerImpl{
		ID:      114,
		TraitID: 18,
		// nolint
		MaleQuestion: `Иногда я говорю глупости ради шутки, чтобы удивить людей и посмотреть, что они на это скажут:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[115] = &KettelItemAnswerImpl{
		ID:      115,
		TraitID: 19,
		// nolint
		MaleQuestion: `Мне бы понравилось быть газетным критиком в разделе драмы, театра, концертов:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[116] = &KettelItemAnswerImpl{
		ID:      116,
		TraitID: 19,
		// nolint
		MaleQuestion: `У меня никогда не бывает потребности что-нибудь рисовать или вертеть в руках, ерзать на месте, когда приходится долго сидеть на собрании:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
	}

	ItemsLib[117] = &KettelItemAnswerImpl{
		ID:      117,
		TraitID: 20,
		// nolint
		MaleQuestion: `Если кто-нибудь говорит мне что-то неправильное, то я скорее подумаю:`,
		Variants: []string{
			`он – лжец`,
			`не уверен`,
			`по-видимому, он плохо информирован`,
		},
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
	}

	ItemsLib[120] = &KettelItemAnswerImpl{
		ID:      120,
		TraitID: 22,
		// nolint
		MaleQuestion: `Торжественность и величие традиционных церемоний следует сохранить:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
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
	}

	ItemsLib[123] = &KettelItemAnswerImpl{
		ID:      123,
		TraitID: 24,
		// nolint
		MaleQuestion: `У меня бывают периоды, когда мне трудно избавиться от чувства жалости к себе:`,
		Variants: []string{
			`часто`,
			`иногда`,
			`никогда`,
		},
	}

	ItemsLib[124] = &KettelItemAnswerImpl{
		ID:      124,
		TraitID: 25,
		// nolint
		MaleQuestion: `Часто я слишком быстро начинаю сердиться на людей:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
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
	}

	ItemsLib[126] = &KettelItemAnswerImpl{
		ID:      126,
		TraitID: 10,
		// nolint
		MaleQuestion: `Если бы зарплата была одинаковой, то я предпочел бы быть:`,
		Variants: []string{
			`адвокатом`,
			`не уверен`,
			`пилотом или капитаном судна`,
		},
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
	}

	ItemsLib[129] = &KettelItemAnswerImpl{
		ID:      129,
		TraitID: 12,
		// nolint
		MaleQuestion: `Когда приходит время осуществить то, что я планировал и на что надеялся, я обнаруживаю, что уже пропало желание делать это:`,
		// nolint
		FemaleQuestion: `Когда приходит время осуществить то, что я планировала и на что надеялась, я обнаруживаю, что уже пропало желание делать это:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
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
	}

	ItemsLib[132] = &KettelItemAnswerImpl{
		ID:      132,
		TraitID: 14,
		// nolint
		MaleQuestion: `Много свободного времени я провожу в разговорах с друзьями о прошлых развлечениях, от которых я получал удовольствие:`,
		// nolint
		FemaleQuestion: `Много свободного времени я провожу в разговорах с подругами и друзьями о прошлых развлечениях, от которых я получала удовольствие:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[133] = &KettelItemAnswerImpl{
		ID:      133,
		TraitID: 14,
		// nolint
		MaleQuestion: `Мне нравится устраивать какие-нибудь смелые рискованные выходки «смеха ради»:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
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
	}

	ItemsLib[137] = &KettelItemAnswerImpl{
		ID:      137,
		TraitID: 17,
		// nolint
		MaleQuestion: `Я люблю музыку:`,
		Variants: []string{
			`легкую, живую`,
			`нечто среднее`,
			`чувствительную`,
		},
	}

	ItemsLib[138] = &KettelItemAnswerImpl{
		ID:      138,
		TraitID: 17,
		// nolint
		MaleQuestion: `Красота поэмы восхищает меня больше, чем красота хорошо сделанного оружия:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[139] = &KettelItemAnswerImpl{
		ID:      139,
		TraitID: 18,
		// nolint
		MaleQuestion: `Если мое удачное замечание остается незамеченным окружающими, то я:`,
		Variants: []string{
			`смирюсь с этим`,
			`нечто среднее`,
			`даю людям возможность услышать его еще раз`,
		},
	}

	ItemsLib[140] = &KettelItemAnswerImpl{
		ID:      140,
		TraitID: 19,
		// nolint
		MaleQuestion: `Мне бы понравилось работать фотокорреспондентом:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[141] = &KettelItemAnswerImpl{
		ID:      141,
		TraitID: 19,
		// nolint
		MaleQuestion: `Нужно быть осторожным в общении с незнакомыми, так как можно, например, заразиться:`,
		// nolint
		FemaleQuestion: `Нужно быть осторожной в общении с незнакомыми, так как можно, например, заразиться:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[142] = &KettelItemAnswerImpl{
		ID:      142,
		TraitID: 20,
		// nolint
		MaleQuestion: `При поездке за границу я бы предпочел быть под руководством экскурсовода, чем самому планировать маршрут:`,
		// nolint
		FemaleQuestion: `При поездке за границу я бы предпочла быть под руководством экскурсовода, чем самой планировать маршрут:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[143] = &KettelItemAnswerImpl{
		ID:      143,
		TraitID: 21,
		// nolint
		MaleQuestion: `Меня справедливо считают упорным и трудолюбивым, но не слишком преуспевающим человеком:`,
		// nolint
		FemaleQuestion: `Меня справедливо считают упорной и трудолюбивой, но не слишком преуспевающей:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[144] = &KettelItemAnswerImpl{
		ID:      144,
		TraitID: 21,
		// nolint
		MaleQuestion: `Если люди пользуются моим хорошим отношением в своих интересах, то я не возмущаюсь этим и вскоре об этом забываю:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
	}

	ItemsLib[145] = &KettelItemAnswerImpl{
		ID:      145,
		TraitID: 22,
		// nolint
		MaleQuestion: `Если при обсуждении какого-либо вопроса среди участников возникает ожесточенный спор, то я предпочитаю:`,
		Variants: []string{
			`увидеть, кто же «победил»`,
			`нечто среднее`,
			`чтобы спор разрешился мирно`,
		},
	}

	ItemsLib[146] = &KettelItemAnswerImpl{
		ID:      146,
		TraitID: 23,
		// nolint
		MaleQuestion: `Я предпочитаю планировать что-либо самостоятельно, без вмешательства и предложений со стороны других:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
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
	}

	ItemsLib[148] = &KettelItemAnswerImpl{
		ID:      148,
		TraitID: 24,
		// nolint
		MaleQuestion: `Я твердо верю, что начальник может быть не всегда прав, но он всегда имеет право быть начальником:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[149] = &KettelItemAnswerImpl{
		ID:      149,
		TraitID: 25,
		// nolint
		MaleQuestion: `Когда я думаю обо всем, что еще предстоит сделать, у меня появляется чувство напряженности:`,
		Variants: []string{
			`да`,
			`иногда`,
			`нет`,
		},
	}

	ItemsLib[150] = &KettelItemAnswerImpl{
		ID:      150,
		TraitID: 25,
		// nolint
		MaleQuestion: `Когда зрители мне что-либо кричат во время игры, меня это не трогает:`,
		Variants: []string{
			`верно`,
			`нечто среднее`,
			`неверно`,
		},
	}

	ItemsLib[151] = &KettelItemAnswerImpl{
		ID:      151,
		TraitID: 10,
		// nolint
		MaleQuestion: `Интереснее быть:`,
		Variants: []string{
			`художником`,
			`не уверен`,
			`организатором культурных развлечений`,
		},
	}

	ItemsLib[152] = &KettelItemAnswerImpl{
		ID:      152,
		TraitID: 11,
		// nolint
		MaleQuestion: `Которое из следующих слов не относится к двум другим:`,
		Variants: []string{
			`любые`,
			`некоторые`,
			`большинство`,
		},
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
	}

	ItemsLib[157] = &KettelItemAnswerImpl{
		ID:      157,
		TraitID: 14,
		// nolint
		MaleQuestion: `Мне больше нравится в одежде спокойная корректность, чем бросающаяся в глаза индивидуальность:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
	}

	ItemsLib[158] = &KettelItemAnswerImpl{
		ID:      158,
		TraitID: 14,
		// nolint
		MaleQuestion: `Мне больше нравится провести вечер за спокойным любимым занятием, чем в оживленной компании:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
	}

	ItemsLib[159] = &KettelItemAnswerImpl{
		ID:      159,
		TraitID: 15,
		// nolint
		MaleQuestion: `Я не обращаю внимания на доброжелательные советы других, даже когда эти советы могли бы быть полезными:`,
		Variants: []string{
			`иногда`,
			`почти никогда`,
			`никогда`,
		},
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
	}

	ItemsLib[163] = &KettelItemAnswerImpl{
		ID:      163,
		TraitID: 17,
		// nolint
		MaleQuestion: `В школе я предпочитал:`,
		// nolint
		FemaleQuestion: `В школе я предпочитала:`,
		Variants: []string{
			`русский язык и литературу`,
			`не уверен`,
			`математику или арифметику`,
		},
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
	}

	ItemsLib[167] = &KettelItemAnswerImpl{
		ID:      167,
		TraitID: 20,
		// nolint
		MaleQuestion: `В воспитании важнее:`,
		Variants: []string{
			`относиться к ребенку с достаточной любовью`,
			`нечто среднее`,
			`выработать нужные привычки и отношение к жизни`,
		},
	}

	ItemsLib[168] = &KettelItemAnswerImpl{
		ID:      168,
		TraitID: 21,
		// nolint
		MaleQuestion: `Люди считают меня положительным, спокойным человеком, которого не трогают превратности судьбы:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[169] = &KettelItemAnswerImpl{
		ID:      169,
		TraitID: 22,
		// nolint
		MaleQuestion: `Я считаю, что общество должно руководствоваться разумом и отбросить старые привычки или ненужные традиции:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[170] = &KettelItemAnswerImpl{
		ID:      170,
		TraitID: 22,
		// nolint
		MaleQuestion: `Думаю, что в современном мире важнее разрешить:`,
		Variants: []string{
			`вопросы нравственности`,
			`не уверен`,
			`разногласия между странами мира`,
		},
	}

	ItemsLib[171] = &KettelItemAnswerImpl{
		ID:      171,
		TraitID: 23,
		// nolint
		MaleQuestion: `Я лучше усваиваю материал:`,
		Variants: []string{
			`читая хорошо написанную книгу`,
			`нечто среднее`,
			`участвуя в обсуждении вопроса`,
		},
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
	}

	ItemsLib[175] = &KettelItemAnswerImpl{
		ID:      175,
		TraitID: 25,
		// nolint
		MaleQuestion: `Под влиянием момента я редко говорю вещи, о которых потом очень сожалею:`,
		Variants: []string{
			`верно`,
			`не уверен`,
			`неверно`,
		},
	}

	ItemsLib[176] = &KettelItemAnswerImpl{
		ID:      176,
		TraitID: 10,
		// nolint
		MaleQuestion: `Если бы меня попросили участвовать в шефской деятельности, то я бы:`,
		Variants: []string{
			`согласился`,
			`не уверен`,
			`вежливо сказал, что занят`,
		},
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
	}

	ItemsLib[180] = &KettelItemAnswerImpl{
		ID:      180,
		TraitID: 13,
		// nolint
		MaleQuestion: `Меня считают человеком, которому обычно в голову приходят хорошие идеи, когда нужно разрешить какую-либо проблему:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
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
	}

	ItemsLib[182] = &KettelItemAnswerImpl{
		ID:      182,
		TraitID: 14,
		// nolint
		MaleQuestion: `Меня считают человеком, полным энтузиазма:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[183] = &KettelItemAnswerImpl{
		ID:      183,
		TraitID: 14,
		// nolint
		MaleQuestion: `Мне нравится работа, которая требует перемен, разнообразия, командировок, даже если она связана с некоторой опасностью:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
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
	}

	ItemsLib[185] = &KettelItemAnswerImpl{
		ID:      185,
		TraitID: 15,
		// nolint
		MaleQuestion: `Мне нравится работа, требующая добросовестного отношения, точных навыков и умений:`,
		Variants: []string{
			`да`,
			`нечто среднее`,
			`нет`,
		},
	}

	ItemsLib[186] = &KettelItemAnswerImpl{
		ID:      186,
		TraitID: 16,
		// nolint
		MaleQuestion: `Я отношусь к типу энергичных людей, которые всегда заняты:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}

	ItemsLib[187] = &KettelItemAnswerImpl{
		ID:      187,
		TraitID: 0,
		// nolint
		MaleQuestion: `Я уверен в том, что не пропустил ни одного вопроса и на все ответил как следует:`,
		// nolint
		FemaleQuestion: `Я уверена в том, что не пропустила ни одного вопроса и на все ответила как следует:`,
		Variants: []string{
			`да`,
			`не уверен`,
			`нет`,
		},
	}
}
