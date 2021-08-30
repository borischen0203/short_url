# **Go URL Shortener using MongoDB**

## Demo Link

 - The app may sleep without using. Just wait for a few seconds to wake it up.

- https://url-demo.herokuapp.com/

<p align="left">
    <img src="https://i.imgur.com/z2YPfa2.png" alt="Sample"  width="623" height="365" >
    <p align="left">
</p>

## Introduce
    This is a prototype of URL Shortener project. You can find the other one that I rewrite
    in my repo. The new version is pure backend part without UI, I used gin framework and
    make code more readable.
    Link: https://github.com/borischen0203/URL-shortener



## Features:

- Generate short URL from long URL
- Redirect to long URL by generated URL
- URL generation supports custom alias


## How to run in local


### Instruction:

      1. Execute main.go file.
      2. Open browser on http://localhots:8000/
      3. Input the long URL in the field[Enter a long URL].
      4. (Optional)Input the the custom ID in the field[Custom alias].
         (make sure you input valid characters(A-Z,a-z,0-9))
      5. Submit
      6. Get the short URL!


<p align="left">
    <img src="https://i.imgur.com/B7Q47kh.png" alt="Sample"  width="426" height="308" >
    <p align="left">
</p>

```
   7. If the page shows the Alias is not available(be used), please try another one!
```

<p align="left">
    <img src="https://i.imgur.com/lbBe18Z.png" alt="Sample"  width="426" height="308" >
</p>

    8. Input the generated short URL in browser, it will redirect to original URL.

## Tech Stack
    - MongoDB
    - HTML
    - CSS
    - Go
    - Heroku deploy
