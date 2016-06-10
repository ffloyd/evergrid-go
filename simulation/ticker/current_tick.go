package ticker

/*
CurrentTick представляет из себя число, которое синхронизируется с Ticker и
всегда содержит актуальное значение tick'а.

В первую очередь информация о наступлении
нового тика поступает в эти сущности, поэтому значение в нем может на единицу опережать
значения подконтрольных Ticker'у Syncable.

Единственный корректный способ инициализоровать эту структуру - это использовать
cоответсвующий метод Ticker'а.
*/
type CurrentTick int

func (ct *CurrentTick) connect(ch chan int) {
	for {
		nextVal := <-ch
		if nextVal == -1 {
			break
		}
		*ct = CurrentTick(nextVal)
	}
}
