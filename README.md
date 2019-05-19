# memo - mempool.observer

<img align="right" width="159px" src="https://raw.githubusercontent.com/0xB10C/memo/v2-master/www/img/brand-icon.png">

[mempool.observer](https://mempool.observer) visualizes various statistics around my Bitcoin memory pool (mempool).
Seemingly stuck and longtime-unconfirmed transactions can be quite annoying for users transacting on the Bitcoin network.
The idea of mempool.observer is to provide users with information about unconfirmed transactions and transaction fees.

> This is v2-master of memo. 
> A full project refresh.

## Project Structure

This repository contains the following files:
- TODO
- TODO

The wiki cointains two architecture overview's: One of the [whole project](https://github.com/0xB10C/memo/wiki/Infrastructure-memo-v2). And one specifically for [memod](https://github.com/0xB10C/memo/wiki/memod-architecture).


## Project History

I've started building the first version of mempool.observer mid 2017 as my first Bitcoin related project.
I was (and still am) motivated by presumably Greg Maxwell's words you can see in the footer.
Later this year the bitcoin transaction fees rose and I had quite some traffic.
The high fees where caused by a huge transaction flood as the price rose to $20k.
I regularly had problems with long running scripts due to querying and processing the huge mempool on a low end VPS.
However due to time constrains I weren't able to work on increasing the performance and the fees were quite low at that time anyway.
This resulted in mempool.observer v1 dieing the not-maintained death sometime in 2018.

I've focused full time on Bitcoin in spring 2019 and spend a part of that time to work on v2.
V2 is a full rewrite of mempool.observer - only the idea, license and the quote from Maxwell remained.
The goal is to offer way more than v1 did, but to build it with performance and maintainability in mind.
I'm open for ideas and feedback.


## Licencse
- TODO
---


>"What's going to happen to Bitcoin?" is the wrong question. The right question is "What are you going to contribute?" &mdash; <cite>[Greg Maxwell](https://github.com/gmaxwell)</cite>
