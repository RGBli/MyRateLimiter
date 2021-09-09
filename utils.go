package myratelimiter

func sliceSum(nums []int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

func setZero(nums []int) {
	for i := 0; i < len(nums); i++ {
		nums[0] = 0
	}
}
