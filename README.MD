# LINE Sticker Downloader
Disclaimer: This is not a project for stealing or reselling or redistributing Line stickers.

-----

LINE (the chat app) has a huge library of stickers that people really like and are unfortunately not available on Whatsapp.

To name a few:
- BT21
- Line characters (Brown & Cony)
- gudetama???
- Some cute rabbits

Long story short, I wanna be that show off with unique stickers on Whatsapp.

## Notes
Over time when the people from LINE finds this, they may change their website to deter scrapping of stickers.<br>
Or perhaps they are aware people are doing so, but do not really care since the exported content cannot be used on LINE anyway.

Regardless, this was never meant to abuse the system, but for personal use only.

## Technical Information
Built using Golang.<br>
Using [apng2gif](https://sourceforge.net/projects/apng2gif/) for converting apng to gif. Great project btw.<br>
Works on Windows only. Buy me a Macbook ~wink wink

2 directory [output, output-gif] will be created if it does not exist.<br>
-> output : apng from the website will be stored here<br>
-> output-gif: apng from output will be converted to gif and stored here. If sticker is a static sticker, it will not be converted<br>
-> bin: the folder where apng2gif binary is stored

### Quick Run of How It Works
1. Download the website and store the HTML as text
2. Regex out a special pattern that the website is using to keep the sticker APNG links.
3. Store that information as a struct/JSON.
4. Using goroutines, loop through that information, download and convert the APNG to GIF.
5. GIF found in output-gif. Now we can manually create Whatsapp stickers yeah...

## FAQ
<strong>"Boohoo why complicate things with Golang when you can do it with Python"</strong><br>
No reasons. I just enjoy the thought of having my exports running concurrently. Also, it will run <strong>at least</strong> 1 second faster than Python.

<strong>"If you used Python, you could support all the OS instead of Windows only"</strong><br>
Wrong. Do your research before assuming. 

<strong>"Is this legal tho"</strong><br>
Hmmmmmmm

<strong>"I have heard web scrapping uses stuff like Selenium"</strong><br>
Yes, Selenium is a great tool that is being used too and is also a viable option.<br>
Fortunately for us, the LINE sticker shop is a server rendered website, meaning a GET request to the main URL will return all the information we need.<br>
This lets us use an approach which relies on one less library and will definitely run faster.

<strong>"Any plans for update/maintenance if LINE update their website?"</strong><br>
Yes only if I see a sticker that I want. Feel free to clone this repo or pull requests and I will definitely take a look.

<strong>"Why though... just use Whatsapp stickers?!"</strong><br>
Because I can therefore I will.

In all seriously, do drop me a message if there are any questions :D
