package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Buku struct {
	KodeBuku string 
	JudulBuku string
	Pengarang string 
	Penerbit string 
	JumlahHalaman int 
	TahunTerbit int 
	
}

var listBuku []Buku


//Fungsi menambahkan data buku
func TambahBuku() {

	inputanUser := bufio.NewReader(os.Stdin)
	judulBuku := ""
	pengarang := ""
	penerbit := ""
	jumlahHalaman := 0
	tahunTerbit := 0
	fmt.Println("=================================")
	fmt.Println("Tambah buku")
	fmt.Println("=================================")
	fmt.Println(listBuku)

	drafBuku := []Buku{}

	for{
		//Kode buku yang sama akan di timpa atau di update
		fmt.Print("Silahkan Masukan kode buku : ")
		kodeBuku, err := inputanUser.ReadString('\n')
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		kodeBuku = strings.Replace(
			kodeBuku,
			"\n",
			"",
			1)
		kodeBuku = strings.Replace(
			kodeBuku,
			"\r",
			"",
			1)

		fmt.Print("Silahkan masukkan judul buku : ")
		_, err = fmt.Scanln(&judulBuku)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		fmt.Print("Silahkan Masukan pengarang : ")
		_, err = fmt.Scanln(&pengarang)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}
		fmt.Print("Silahkan Masukan penerbit : ")
		_, err = fmt.Scanln(&penerbit)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}
		fmt.Print("Silahkan Masukan jumlah halaman : ")
		_, err = fmt.Scanln(&jumlahHalaman)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}
		fmt.Print("Silahkan Masukan tahun terbit : ")
		_, err = fmt.Scanln(&tahunTerbit)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}
		
		drafBuku = append(drafBuku, Buku{
		KodeBuku : fmt.Sprintf("book-%s", kodeBuku),
		JudulBuku : judulBuku,
		Pengarang : pengarang,
		Penerbit :penerbit ,
		JumlahHalaman :jumlahHalaman ,
		TahunTerbit : tahunTerbit,
		})

		pilihanmenu := 0 
		fmt.Print("ketik 1 untuk tambah pesanan lagi, ketik 0 untuk keluar :")
		_, err = fmt.Scanln(&pilihanmenu)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		if pilihanmenu == 0 {
			break
		}

}
	fmt.Println("Menambah Buku...")

	_ = os.Mkdir("books", 0777)

	ch := make(chan Buku)

	wg := sync.WaitGroup{}

	jumlahPustakawan := 5

	for i := 0; i < jumlahPustakawan; i++ {
		wg.Add(1)
		go simpanBuku(ch, &wg, i)
	}

	// Mengirimkan data ke channel
	for _, books := range drafBuku {
		ch <- books
	}
	
	close(ch)

	wg.Wait()
	fmt.Println("Berhasil Menambah data buku!")
}

func simpanBuku(ch <-chan Buku, wg *sync.WaitGroup, noPustakawan int) {

	for bukuDisimpan := range ch {
		dataJson, err := json.Marshal(&bukuDisimpan)

		if err != nil {
			fmt.Println("Terjadi error:", err)
		}

		err = os.WriteFile(fmt.Sprintf("books/%s.json", bukuDisimpan.KodeBuku), dataJson, 0644)
		if err != nil {
			fmt.Println("Terjadi error:", err)
		}

		fmt.Printf("Pelayan No %d Memproses Buku dengan kodeBuku : %s!\n", noPustakawan, bukuDisimpan.KodeBuku)
	}
	wg.Done()
}

func lihatBuku(ch <-chan string, chBuku chan Buku, wg *sync.WaitGroup) {
	var books Buku
	for idBuku := range ch {
		dataJSON, err := os.ReadFile(fmt.Sprintf("books/%s", idBuku))
		if err != nil {
			fmt.Println("Terjadi error:", err)
		}

		err = json.Unmarshal(dataJSON, &books)
		if err != nil {
			fmt.Println("Terjadi error:", err)
		}

		chBuku <- books
	}
	wg.Done()
}
//Fungsi melihat data buku
func LihatBuku() {
	fmt.Println("=================================")
	fmt.Println("Lihat Pesanan")
	fmt.Println("=================================")
	fmt.Println("Memuat data ...")
	listBuku = []Buku{}

	listJsonBuku, err := os.ReadDir("buku")
	if err != nil {
		fmt.Println("Terjadi error: ", err)
	}

	wg := sync.WaitGroup{}

	ch := make(chan string)
	chBuku := make(chan Buku, len(listJsonBuku))

	jumlahPustakawan := 5

	for i := 0; i < jumlahPustakawan; i++ {
		wg.Add(1)
		go lihatBuku(ch, chBuku, &wg)
	}

	for _, fileBuku := range listJsonBuku {
		ch <- fileBuku.Name()
	}

	close(ch)

	wg.Wait()

	close(chBuku)

	for dataBuku := range chBuku {
		listBuku = append(listBuku, dataBuku)
	}

	if len(listBuku) < 1 {
		fmt.Println("Buku tidak ada")
	}

	for urutan, books := range listBuku {
		fmt.Printf("%d. kode buku : %s\n   judul buku : %s\n   pengarang buku : %s\n   penerbit buku: %s\n   jumlah halaman buku: %d\n   tahun terbit buku : %d\n",
			urutan+1,
			books.KodeBuku,
			books.JudulBuku,
			books.Pengarang,
			books.Penerbit,
			books.JumlahHalaman,
			books.TahunTerbit,
		)
	}

}



func DetailBuku(kode string) {
	fmt.Println("\r")
	fmt.Println("======")
	fmt.Println("Detail Buku")
	fmt.Println("======")

	var isBook bool

	for _, book := range listBuku {
		if book.KodeBuku == kode {
			isBook = true
			fmt.Printf("Kode Buku : %s\n", book.KodeBuku)
			fmt.Printf("Judul Buku : %s\n", book.JudulBuku)
			fmt.Printf("Pengarang Buku : %s\n", book.Pengarang)
			fmt.Printf("Penerbit Buku : %s\n", book.Penerbit)
			fmt.Printf("Jumlah Halaman : %d\n", book.JumlahHalaman)
			fmt.Printf("Tahun Terbit : %d\n", book.TahunTerbit)
			break
		}
	}

	if !isBook {
		fmt.Println("Kode Buku Salah Atau Tidak Ada")
	}
}

func updateBuku(kode string) {
	fmt.Println("\r")
	DetailBuku(kode)

	fmt.Println("======")
	fmt.Println("Edit Buku")
	fmt.Println("======")

	var book Buku

	fmt.Print("Masukkan Code Buku : ")
	_, err := fmt.Scanln(&book.KodeBuku)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukkan Judul Buku : ")
	_, err = fmt.Scanln(&book.JudulBuku)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukkan Pengarang Buku : ")
	_, err = fmt.Scanln(&book.Pengarang)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukkan Penerbit Buku : ")
	_, err = fmt.Scanln(&book.Penerbit)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukkan Total Halaman : ")
	_, err = fmt.Scanln(&book.JumlahHalaman)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukkan Tahun Terbit : ")
	_, err = fmt.Scanln(&book.TahunTerbit)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Println(book)

	for i, b := range listBuku {
		if b.KodeBuku == kode {
			listBuku[i] = book
			break
		}
	}
}


//fungsi menghapus data buku
func HapusBuku(kode string) {
	fmt.Println("=================================")
	fmt.Println("Hapus data buku")
	fmt.Println("=================================")
	LihatBuku()
	fmt.Println("=================================")
	var isBook bool
	for i, book := range  listBuku{
		if book.KodeBuku == kode {
			isBook = true
			err := os.Remove(fmt.Sprintf("buku/%s.json", listBuku[i].KodeBuku))
			if err != nil {
				fmt.Println("Terjadi error:", err)
			}
			fmt.Println("Buku Berhasil Dihapus")
			break
		}
	}
	if !isBook {
		fmt.Println("Kode Buku Salah Atau Tidak Ada")
	}
}


//fungsi utama
func main() {
	var pilihanMenu int
	fmt.Println("Sistem Manajemen Manajemen Daftar Buku Perpustakaan")
	fmt.Println("=================================")
	fmt.Println("Silahkan Inputkan pilihan anda : ")
	fmt.Println("1. Tambah Buku baru")
	fmt.Println("2. Lihat Buku")
	fmt.Println("3. Lihat Detail Buku")
	fmt.Println("4. Edit Buku")
	fmt.Println("5. Hapus Buku")
	fmt.Println("6. Keluar")
	fmt.Println("=================================")
	fmt.Print("Inputkan Pilihan disini : ")
	_, err := fmt.Scanln(&pilihanMenu)
	if err != nil {
		fmt.Println("Terjadi error:", err)
	}

	switch pilihanMenu {
	case 1:
		TambahBuku()
	case 2:
		LihatBuku()
	case 3:
		var pilihDetail string
		LihatBuku()
		fmt.Print("Masukkan Kode Buku : ")
		_, err := fmt.Scanln(&pilihDetail)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}
		DetailBuku(pilihDetail)
		
	case 4:
		var pilihanUpdate string
		LihatBuku()
		fmt.Print("Masukkan Kode Buku Yang Akan DiEdit : ")
		_, err := fmt.Scanln(&pilihanUpdate)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}
		updateBuku(pilihanUpdate)
		
	case 5:
		var pilihanHapus string
		LihatBuku()
		fmt.Print("Masukkan Kode Buku Yang Akan Dihapus : ")
		_, err := fmt.Scanln(&pilihanHapus)
		if err != nil {
			fmt.Println("Terjadi Kesalahan : ", err)
			return
		}
		HapusBuku(pilihanHapus)
	case 6:
		os.Exit(0)
	}

	main()
}