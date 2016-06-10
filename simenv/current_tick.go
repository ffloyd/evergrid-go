package simenv

/*
CurrentTick представляет из себя число, которое синхронизируется с SimEnv и
всегда содержит актуальное значение tick'а.

В первую очередь информация о наступлении
нового тика поступает в эти сущности, поэтому значение в нем может на единицу опережать
значения подконтрольных SimEnv'у Agent'ов.

Единственный корректный способ инициализоровать эту структуру - это использовать
cоответсвующий метод SimEnv'а.
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
