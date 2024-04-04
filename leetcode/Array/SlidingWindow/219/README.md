### Contains Duplicate II


# Given an integer array nums and an integer k, return true if there are two distinct indices i and j in the array such that nums[i] == nums[j] and abs(i - j) <= k.

Example 1:

Input: nums = [1,2,3,1], k = 3
Output: true

Example 2:

Input: nums = [1,0,1,1], k = 1
Output: true

Example 3:

Input: nums = [1,2,3,1,2,3], k = 2
Output: false

# Đề bài yêu cầu bạn kiểm tra xem trong một mảng số nguyên nums có tồn tại hai chỉ số khác nhau i và j sao cho nums[i] bằng nums[j] và khoảng cách tuyệt đối giữa i và j không vượt quá k.

Ví dụ 1:

Input: nums = [1,2,3,1], k = 3
Output: true
Trong mảng nums, số 1 xuất hiện ở chỉ số 0 và 3, với khoảng cách là 3 - 0 = 3, thỏa mãn điều kiện.

Ví dụ 2:

Input: nums = [1,0,1,1], k = 1
Output: true
Trong mảng nums, số 1 xuất hiện ở chỉ số 0, 2, và 3. Chúng ta có cặp (0, 2) và (2, 3) thỏa mãn điều kiện với khoảng cách tối đa là 1.

Ví dụ 3:

Input: nums = [1,2,3,1,2,3], k = 2
Output: false
Trong mảng nums, không có cặp nào thỏa mãn điều kiện với khoảng cách tối đa là 2. Số 1 xuất hiện ở chỉ số 0 và 3 nhưng khoảng cách là 3 - 0 = 3 > 2.