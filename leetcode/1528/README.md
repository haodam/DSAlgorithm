1528. Shuffle String
You are given a string s and an integer array indices of the same length. 
The string s will be shuffled such that the character at the ith position moves to indices[i] in the shuffled string.
Return the shuffled string

Input: s = "codeleet", indices = [4,5,6,7,0,2,1,3]
Output: "leetcode"
Explanation: As shown, "codeleet" becomes "leetcode" after shuffling.

Example 2:

Input: s = "abc", indices = [0,1,2]
Output: "abc"
Explanation: After shuffling, each character remains in its position.


- Ví dụ 1: cho một chuỗi s = "codeleet" và một mảng indices = [4,5,6,7,0,2,1,3] ứng với các giá trị trong mảng, 
ví dụ "c" sẽ ứng với vị trí thứ 4 trong mảng indices , nhiệm vụ là hãy sắp xếp lại theo đúng thứ tự từ lớn 
đến bé indices = [0,1,2,3,4,5,6,7] thì chuỗi s = "leetcode".

- Thuật giải: 
  + Tạo một mảng byte bằng với độ dài mảng indices.
  + Chạy vòng lặp for có index và value.
  + Biến result[value] = s[index] ví dụ chữ c có value = [4] sẽ truy xuất lấy chữ c có index = 0 ở chuỗi s
  đặt vào vào vị trí thứ 4 , index = 0, value = 4 , tương tự chữ o index = 1 , value = 5.   
  + kết quả trả về sẽ ép kiểu về thành chuỗi string.

- Độ phức tạp của thuật toán:
  + Độ phức tạp về thời gian là: 0(n)
  + Độ phức tạp về không gian lưu chữ là: 0(n) 



