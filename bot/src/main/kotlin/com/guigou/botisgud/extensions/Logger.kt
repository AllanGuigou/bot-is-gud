package com.guigou.botisgud.extensions

import org.slf4j.Logger
import org.slf4j.LoggerFactory

inline fun <reified R : Any> R.logger(): Logger = LoggerFactory.getLogger(R::class.java.declaringClass)
