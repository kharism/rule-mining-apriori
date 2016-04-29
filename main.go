package main
import(
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"sort"
	. "github.com/kharism/hashtree"
	"os"
	//"encoding/csv"
)
func contains(arr []int,i int)bool{
	for _,j:=range(arr){
		if j==i{
			return true
		}
	}
	return false
}
func isSame(a,b []int)bool{
	if len(a)!=len(b){
		return false
	}
	for i,_:=range(a){
		if a[i]!=b[i]{
			return false
		}
	}
	return true
}
type HashCount struct{
	
} 
func main(){
	threshold:=0.00005
	itemSetLen := 3
	file,_:=ioutil.ReadFile("5000-out1.csv")
	data := strings.Split(string(file),"\n")
	//root := make(map[int]*HashNode)
	recordCount:=len(data)-1
	fmt.Println(len(data))
	itemCount:=make(map[int]int)
	itemCountF:=make(map[int]float64)
	transactions:=[][]int{}
	for _,i:=range(data){
		if len(i)==1{
			continue
		}
		i=strings.TrimRight(i,", ")
		splitted := strings.Split(i,",")
		
		transaction:=[]int{}
		for j:=1;j<len(splitted);j++{
			itemI,err:=strconv.Atoi(strings.TrimRight(strings.TrimLeft(splitted[j]," ")," "))
			
			if err!=nil{
				fmt.Println(err.Error())
				break
			} else{
				transaction = append(transaction,itemI)
				itemCount[itemI]+=1
			}
			
		}
		
		sort.Ints(transaction)
		if len(transaction)==0{
			fmt.Println(i)
			continue
		}
		//fmt.Println("Adding transaction",transaction)
		//_,ok:=root[transaction[0]]
		//if !ok{
		//	root[transaction[0]] = &HashNode{Key:transaction[0]}
		//}
		//root[transaction[0]].AddValueRecursive(transaction[1:])
		transactions = append(transactions,transaction)		
	}
	keys := []int{}
	frequentItemset:=[][]int{}
	subFrequentItemset:=[][]int{}
	frequentItems:=[]int{}
	for key:=range(itemCount){
		itemCountF[key] = float64(itemCount[key])/float64(recordCount)
		
		keys = append(keys,key)
	}
	sort.Ints(keys)
	for _,key:=range(keys){
		if itemCountF[key]>=threshold{
			frequentItems = append(frequentItems,key)
			frequentItemset = append(frequentItemset,[]int{key})
		}
	}
	//generating FrequentItemset for n>=2
	root:=&HashNode{}
	subRoot:=&HashNode{}
	
	for count:=1;count<itemSetLen;count++{
		newFrequentItemset:=[][]int{}
		// i = array of int 
		for _,i:=range(frequentItemset){
			curFreq :=1.0
			for _,j:=range(i){
				curFreq=curFreq*itemCountF[j]
			}
			for _,j:=range(frequentItems){
				if i[len(i)-1]>=j{
					continue
				}
				if curFreq*itemCountF[j]>=threshold{
					k:=append(i,j)
					if len(i)==3{
						//fmt.Println(i,j)
					}
					
					//if len(newFrequentItemset)>0 && !isSame(k,newFrequentItemset[len(newFrequentItemset)-1]){
						newFrequentItemset = append(newFrequentItemset,k)
					//}
					
				}
			}
		}
		//add subroot for calculating confidence 
		if count==itemSetLen-2{
			fmt.Println("add subroot")
			fmt.Println("aaa",len(newFrequentItemset))
			subFrequentItemset = newFrequentItemset
			for _,itemset:=range(subFrequentItemset){
				//fmt.Println("adding node ",itemset,"to ")
				subRoot.AddNodeWithoutValue(itemset)
			}
			for _,itemset:=range(transactions){
				if len(itemset)<itemSetLen-1{
					continue
				}
				fmt.Println("Ordered Combination of ",itemset)
				if len(itemset)==itemSetLen-1{
					fmt.Println("Adding value to node ",itemset)
					subRoot.AddValueWithoutCreate(itemset)
				}else{
					fmt.Println("Check Every Possible combination")
					for i:=0;i<=len(itemset)-itemSetLen+1;i++{
						fmt.Println(itemset[i])
						for j:=i+1;j<=len(itemset)-itemSetLen+2;j++{
							k:=append([]int{itemset[i]},itemset[j:itemSetLen+j-2]...)
							fmt.Println("Adding value to node ",k)
							subRoot.AddValueWithoutCreate(k)
							
						}
					}
				}
				
			}
			
		}
		//fmt.Println(newFrequentItemset)
		frequentItemset= newFrequentItemset
	}
	//fmt.Println(frequentItemset)
	//Generating HashTree from frequent itemset
	//5 & 3->3
	//5 & 4->2
	//4 & 5->1
	fmt.Println("aab",len(frequentItemset))
	for _,itemset:=range(frequentItemset){
		root.AddNodeWithoutValue(itemset)
	}
	fmt.Println("Add transactions to root")
	for _,itemset:=range(transactions){
		if len(itemset)<itemSetLen{
			continue
		}
		if len(itemset)==itemSetLen{
			fmt.Println("Add Langsung",itemset)
			root.AddValueRecursive(itemset)
		} else{
			fmt.Println("Ordered Combination of ",itemset)
			for i:=0;i<=len(itemset)-itemSetLen+1;i++{
				//fmt.Println(itemset[i])
				for j:=i+1;j<=len(itemset)-itemSetLen+1;j++{
					k:=append([]int{itemset[i]},itemset[j:itemSetLen+j-1]...)
					fmt.Println(k)
					root.AddValueRecursive(k)
				}
				
			}
		}
		
	}
	fmt.Println("Item Count")
	for _,itemset:=range(frequentItemset){
		j,_:=root.GetValueRecursive(itemset)
		if j>0{
			fmt.Println(itemset,j)
		}
		
	}
	fmt.Println("sub-Item Count")
	for _,itemset:=range(subFrequentItemset){
		j,_:=subRoot.GetValueRecursive(itemset)
		if j>10{
			fmt.Println(itemset,j)
		}
	}
	fmt.Println("Rule Mining")
	fmt.Println("=====")
	for _,itemset:=range(frequentItemset){
		allCount,_:=root.GetValueRecursive(itemset)
		if allCount==0{
			//fmt.Println("Kosong")
			continue
		}
		//fmt.Println("Periksa sub-Itemset",itemset)
		for i:=0;i<=itemSetLen-1;i++{
			//fmt.Println("check appending",i,itemset[i])
			//fmt.Println("check appending2",i+1,itemSetLen-1)
			for j:=i+1;j<itemSetLen-1;j++{
				k:=[]int{itemset[i]}
				//fmt.Println(k,itemset[j:j+itemSetLen-2])
				k = append(k,itemset[j:j+itemSetLen-2]...)
				//fmt.Println(k)
				ua,err:=subRoot.GetValueRecursive(k)
				
				if err==nil && ua>0{
					confidentLevel:=float64(allCount)/float64(ua)
					if ua>=10 && allCount>=10{
						fmt.Println("Confident Level",k,itemset,ua,allCount,confidentLevel)
					}
					
				} else if err!=nil{
					fmt.Println("Gagal Ambil ",k)
					fmt.Println(err.Error())
					os.Exit(1)
				}
				
			}
			//fmt.Println("Finish appending",itemset)
		}
		
	}
	//fmt.Println(root[4].Childs)
	//fmt.Println(root[4].GetValueRecursive([]int{7}))
	
}
