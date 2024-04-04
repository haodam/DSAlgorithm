### 594. Longest Harmonious Subsequence

We define a harmonious array as an array where the difference between its maximum value and its minimum value is exactly 1.
Given an integer array nums, return the length of its longest harmonious subsequence among all its possible subsequences.
A subsequence of array is a sequence that can be derived from the array by deleting some or no elements without changing the order of the remaining elements

Example 1:

Input: nums = [1,3,2,2,5,2,3,7]
Output: 5
Explanation: The longest harmonious subsequence is [3,2,2,2,3].

Example 2:

Input: nums = [1,2,3,4]
Output: 2

Example 3:

Input: nums = [1,1,1,1]
Output: 0


# Bài toán:
Yêu cầu tìm độ dài của dãy con dài nhất trong mảng, sao cho độ chênh lệch giữa giá trị lớn nhất và giá trị
nhỏ nhất của dãy con đó bằng 1. Đây là một số điều cần lưu ý:
Harmonious array: Một mảng được gọi là "harmonious" nếu chênh lệch giữa giá trị lớn nhất và giá trị nhỏ nhất trong mảng đó là đúng 1.
Subsequence: Một dãy con là một dãy mà có thể thu được từ mảng gốc bằng cách xóa đi một số phần tử hoặc không xóa phần
tử nào, mà vẫn giữ nguyên thứ tự của các phần tử còn lại.

# Ví dụ:

Trong ví dụ 1, mảng nums = [1,3,2,2,5,2,3,7], dãy con dài nhất có thể thu được là [3,2,2,2,3], với độ dài là 5.
Trong ví dụ 2, mảng nums = [1,2,3,4], dãy con dài nhất có thể thu được là [1,2] hoặc [2,3], với độ dài là 2.
Trong ví dụ 3, mảng nums = [1,1,1,1], không có dãy con nào thỏa mãn điều kiện, vì không có giá trị nào khác biệt với nhau.