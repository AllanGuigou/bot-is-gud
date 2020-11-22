package com.guigou.botisgud.services

import java.io.File

class WordServiceImpl : WordService {
    // TODO: readLines should not be used for large files
    private val words = File("words").readLines()

    override suspend fun random(): String {
        return words.random()
    }
}
