1502. Can Make Arithmetic Progression From Sequence
- A sequence of numbers is called an arithmetic progression if the difference between any two consecutive elements is the same.
Given an array of numbers arr, return true if the array can be rearranged to form an arithmetic progression. Otherwise, return false.

Example 1:

Input: arr = [3,5,1]
Output: true
Explanation: We can reorder the elements as [1,3,5] or [5,3,1] with differences 2 and -2 respectively, between each consecutive elements.


- Lời giải: Cho một mảng arr[] hãy kiểm tra xem mảng arr có phải là một dãy cấp số cộng không?

- Ví dụ 1: arr = [3,5,1] là một dãy cấp số cộng kêt quả trả về true , ngược lại là false

- Thuật giải: 
    + Tìm 2 giá trị min và max trong mảng arr
    + Tính khoảng cách 2 giá trị cách nhau (diff) , ví dụ 1 và 3 có khoảng cách nhảy là 2
      * diff = 1 --> 1 2 3 4 5 6 7 8 9 | N = 9
      * diff = 2 --> 1 3 5 7 9 | N = 5
    + Công thức tính diff = (max - min ) / (N -1) trong đó N là số phần tử trong mảng
    + Kiểm tra khoảng cách nếu là số chẵn = > true, là số lẻ => false

- Độ phức tạp của thuật toán:
    + Độ phức tạp về thời gian là: 0(n)
    + Độ phức tạp về không gian lưu chữ là: 0(n) 
