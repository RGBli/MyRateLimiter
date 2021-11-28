package myratelimiter

// sliceSum return the sum of an int type slice
func sliceSum(nums []int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

// setSliceZero sets every element in an int type to zero
func setSliceZero(nums []int) {
	for i := 0; i < len(nums); i++ {
		nums[0] = 0
	}
}
