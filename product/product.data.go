package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("Begin load data...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d products loaded...\n", len(productMap.m))
}

func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	productsJson, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)
	err = json.Unmarshal(productsJson, &productList)
	if err != nil {
		log.Fatal(err)
	}

	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}
	return prodMap, nil
}

func getProductIds() []int {
	productMap.RLock()
	productIds := []int{}

	for key := range productMap.m {
		productIds = append(productIds, key)
	}

	productMap.RUnlock()
	sort.Ints(productIds)
	return productIds
}

func getNextProductID() int {
	productIds := getProductIds()
	return productIds[len(productIds)-1] + 1
}

func getProductList() []Product {
	productMap.RLock()
	products := make([]Product, 0, len(productMap.m))

	for _, product := range productMap.m {
		products = append(products, product)
	}
	productMap.RUnlock()
	return products
}

func getProduct(productID int) *Product {
	productMap.RLock()
	defer productMap.RUnlock()

	if product, ok := productMap.m[productID]; ok {
		return &product
	}
	return nil
}

func addOrUpdateProduct(product Product) (int, error) {
	addOrUpdateID := -1

	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)
		if oldProduct == nil {
			return -1, fmt.Errorf("product id [%d] doesn't exist", product.ProductID)
		}
		addOrUpdateID = product.ProductID
	} else {
		addOrUpdateID = getNextProductID()
		product.ProductID = addOrUpdateID
	}
	productMap.Lock()
	productMap.m[addOrUpdateID] = product
	productMap.Unlock()
	return addOrUpdateID, nil
}

func removeProduct(ProductID int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, ProductID)
}

func getTopTenProducts() []Product {
	products := getProductList()
	sort.SliceStable(products, func (i, j int) bool {
		return products[i].QuantityOnHand > products[j].QuantityOnHand
	})
	if len(products) > 10 {
		return products[:10]
	}
	return products
}