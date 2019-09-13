(ns ww.core
  (:gen-class)
  (:require [clojure.string :as str]))


(defn intersperse [v xs]
  (if (< (count xs) 2)
    xs
    (concat [(first xs) v] (intersperse v (rest xs)))))


(defn line-length [words]
  (apply + (intersperse 1 (map count words))))


(defn word-wrap [maxlen words]
  (loop [acc      []
         linebuf  []
         xs       words]
    (if (empty? xs)
      (if (empty? linebuf) acc (conj acc linebuf))
      (let [linebuf' (conj linebuf (first xs))]
        (if (> (line-length linebuf') maxlen)
          (recur (conj acc linebuf)
                 [(first xs)]
                 (rest xs))
          (recur acc
                 linebuf'
                 (rest xs)))))))


(defn words [s]
  (str/split s #"\s+"))


(defn unwords [wrds]
  (apply str (intersperse " " wrds)))


(defn unlines [lines]
  (apply str (intersperse "\n" lines)))


(defn wrap-str [maxlen s]
  (unlines (map unwords (word-wrap maxlen (words s)))))


(defn -main
  [& args]
  (println (wrap-str 80 (slurp *in*))))
