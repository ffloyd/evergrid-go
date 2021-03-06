/*
evergrid-go - это реализация инфраструктурных компонентов для системы распредленный вычислений Evergrid.

Идея этой команды в том, что должна быть одна кодовая база для всех компонентов,
а их запуск должен различаться только вызовом команды. Для реализации CLI используется
пакет github.com/spf13/cobra

На данный момент реализовано только два компонента:

gendata - генерирует сценарии для симуляции

simulator - запускает симуляции

CLI предоставляет краткую справку по возможным командам и опциям их вызова
*/
package main

import "github.com/ffloyd/evergrid-go/cmd"

func main() {
	cmd.Execute()
}
