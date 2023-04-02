# Command-line-filemanager
___________________________________

## Утилита командной строки, позволяющая работать с файлами при помощи команд **wc** и **cat** и флагов.

Утилита написана на Go, её функционал реализован при помощи горутин, одной управляющей горутины, каналов и конвейеров.

Утилита состоит из двух частей: **cat**, **wc**

### cat
**cat** - позволяет читатать файлы и записывать их содержимое в стандартный поток вывода: cat **[OPTION]..** **[FILENAME]..**
* cat (default) - вывести в стандартны поток вывода содержимое файла
* cat -b - пронумеровать непустые выходные строки
* cat -n - пронумеровать все выходные строки

## Пример использования -cat **[OPTION]..** **[FILENAME]..**:
![result1](https://github.com/ellofae/Command-line-filemanager/blob/main/img/Screenshot%20from%202023-04-02%2017-48-56.png?raw=true)
![result1](https://github.com/ellofae/Command-line-filemanager/blob/main/img/Screenshot%20from%202023-04-02%2017-49-07.png?raw=true)

### wc
**wc** позволяет печатает количество новых строк, слов и байт файла при помощи команды: wc **[FILENAME]..**

## Пример использования wc:
![result1](https://github.com/ellofae/Command-line-filemanager/blob/main/img/Screenshot%20from%202023-04-02%2017-48-07.png?raw=true)
