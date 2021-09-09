package myratelimiter

// sliceSum 用于计算切片的和
func sliceSum(nums []int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

// setZero 用于将切片置零
func setZero(nums []int) {
	for i := 0; i < len(nums); i++ {
		nums[0] = 0
	}
}
