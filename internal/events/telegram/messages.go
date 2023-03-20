package telegram

const msgHelp = `Send me any of the genres on the list
Genres:
--------------------------------
		"Короткометражка"
		"Биография"     
		"Боевик"     
		"Вестерн"       
		"Военный"       
		"Детектив"       
		"Для взрослых"  
		"Документальный"
		"Драма"
		"Игра"          
		"История"       
		"Комедия"       
		"Криминал"      
		"Мелодрама"    
		"Музыка"      
		"Мюзикл"       
		"Новости"       
		"Приключения"
		"Семейный"    
		"Спорт"         
		"Триллер"      
		"Ужасы"        
		"Фантастика"   
		"Фэнтези"
--------------------------------

/stat - Outputs statistics on your requests`

const msgHello = "Hi there! 👾\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgNoHistory      = "You don't have stories yet🥺"
	msgAllMovie       = "You have watched all the movies we could show in this genre. " +
		"We apologize😔. Bot in development🤖. We'll fix it soon"
)
