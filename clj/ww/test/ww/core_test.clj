(ns ww.core-test
  (:require [clojure.test :refer :all]
            [ww.core :refer :all]))


(deftest intersperse-test
  (testing "handles edge cases"
    (is (= [] (intersperse 17 [])))
    (is (= [1] (intersperse 17 [1]))))
  (testing "intersperses"
    (is (= [1 17 2 17 3] (intersperse 17 [1 2 3])))))


(deftest line-length-test
  (testing "no words"
    (is (= 0 (line-length []))))
  (testing "one word"
    (is (= 5 (line-length ["hello"]))))
  (testing "many words"
    (is (= 12 (line-length ["hi" "there" "sir"])))))


(deftest word-wrap-test
  (testing "doesn't wrap"
    (is (= [["hi" "there"]] (word-wrap 100 ["hi" "there"]))))
  (testing "wraps"
    (is (= [["hey" "there"] ["you"] ["handsome"]] (word-wrap 10 ["hey" "there" "you" "handsome"])))))


(deftest words-test
  (testing "no breaks"
    (is (= ["hi"] (words "hi"))))
  (testing "breaks"
    (is (= ["a" "b" "c" "d" "e" "f" "g"] (words "a b  c\nd\ne\tf\t   g")))))


(deftest unwords-test
  (testing "unwords"
    (is (= "hi there" (unwords ["hi" "there"])))))


(deftest unlines-test
  (testing "unlines"
    (is (= "hi\nthere" (unlines ["hi" "there"])))))


(deftest wrap-str-test
  (testing "wraps string"
    (is (= "hey there\nyou\nhandsome\nperson"
           (wrap-str 10 "hey there you handsome person")))))
