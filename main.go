package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func bssel() {
	bs := make(chan string, 2) // буффер канал.
	bs <- "первую"             // передал в канал элемент.

	select {
	case act := <-bs: // если прочитать операция будет не блокирующая
		fmt.Println("читать", act)
	case bs <- "секунду": // если буфер заполнен операция заблокируется
		fmt.Println("писать", <-bs, <-bs) //запись в канал
	}
}

func main() {
	unch := make(chan int) // не буфф канал

	go func() {
		time.Sleep(time.Second)
		unch <- 1 // передал значения в канал
	}()

	select {
	// case bs <- "ttt": //запись в буфферизированный канал, если сработает то пишет:
	// 	fmt.Println("Не блокированная запись.")
	case val := <-unch: // если неполучается прочитать то идем дальше
		fmt.Println("Блокировка чтения", val)
	case <-time.After(time.Millisecond * 500): // находит новые данные если не нашел за 0,5s:
		fmt.Println("Время истекло")
	default: // если ничего не можем выполняется:
		fmt.Println("Отключение кейса.")
	}

	res := make(chan int)
	time1 := time.After(time.Second) // передал тайм в канал

	go func() { // функция
		defer close(res) // закрыть канал

		for i := 1; i < 1000; i++ { // цикл от 1 до 1000

			select {
			case <-time1:
				fmt.Println("Время истекло.")
				return
			default:
				time.Sleep(time.Nanosecond) // через наносекунду код закрывается = 66
				res <- i                    // передача цикла в канал

			}
		}

	}()

	for v := range res {
		fmt.Println(v) // то  что получилось в канале
	}

	bssel()
	gr()
}

// все это на практике:

// Select. Graceful shutdown.
func gr() {
	sg := make(chan os.Signal, 1) // буфферизированный канал с сигналами

	signal.Notify(sg, syscall.SIGINT, syscall.SIGTERM) // сигналы

	timer := time.After(10 * time.Second) // запуск таймера на 10 сек

	for {
		select {
		case <-timer: // передали время в канал по истечению которого если не найдется:
			fmt.Println("Время истекло.")
			return
		case sig := <-sg: // попытка получить значения из канала которого приходит сигнал
			fmt.Println("Остановлено сигналом:", sig)
			return
		}
	}
}

/* когда программа запустится и если она услышит syscall.SIGINT, syscall.SIGTERM то
в таком случае в канал отправляются данные какой сигнал получили и остановку програы, все это за 10 секунд поиска */
