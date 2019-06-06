# memo - mempool.observer

[![Build Status](https://travis-ci.com/0xB10C/memo.svg?branch=v2-master)](https://travis-ci.com/0xB10C/memo)

<img align="right" width="280" src="https://raw.githubusercontent.com/0xB10C/memo/v2-master/www/img/brand-icon.png">

[mempool.observer](https://mempool.observer) visualizes various statistics around my Bitcoin memory pool (mempool).
Seemingly stuck and longtime-unconfirmed transactions can be quite annoying for users transacting on the Bitcoin network.
The idea of mempool.observer is to provide users with information about unconfirmed transactions and transaction fees.

> This is v2-master of memo. 
> A full project refresh.

## Project Structure

Folder Structure
```
memo/
├── api/          # Go module functioning as an API returning JSON
├── database/     # Database creation scripts
├── memod/        # Go module functioning as a worker deamon wirting data to database
└── www/          # Statically served HTML, JS and CSS files
```

The wiki cointains two architecture overview's: One of the [whole project](https://github.com/0xB10C/memo/wiki/Infrastructure-memo-v2). And one specifically for [memod](https://github.com/0xB10C/memo/wiki/memod-architecture).


## Project History

I've started building the first version of mempool.observer mid 2017 as my first Bitcoin related project.
I was (and still am) motivated by presumably Greg Maxwell's words:

>"What's going to happen to Bitcoin?" is the wrong question. The right question is "What are you going to contribute?" &mdash; <cite>[Greg Maxwell](https://github.com/gmaxwell)</cite>

Later this year the bitcoin transaction fees rose and I had quite some traffic.
The high fees were caused by a huge transaction flood as the price rose to $20k.
I regularly had problems with long running scripts due to querying and processing the huge mempool on a low end VPS.
Due to time constrains I wasn't able to improve the performance.
This resulted in mempool.observer v1 dieing the not-maintained-death sometime in 2018.

I've focused full time on Bitcoin in spring 2019 and spend a part of that time to work on v2.
V2 is a full rewrite of mempool.observer - only the idea, license and the quote from Maxwell remained.
The goal is to offer way more than v1 did, but build on a foundation with performance and maintainability in mind.
I'm open for ideas and feedback.


## Licencse
This project and all it's files are licensed under a GNU Affero General Public License.

---

>"What's going to happen to Bitcoin?" is the wrong question. The right question is "What are you going to contribute?" &mdash; <cite>[Greg Maxwell](https://github.com/gmaxwell)</cite>

