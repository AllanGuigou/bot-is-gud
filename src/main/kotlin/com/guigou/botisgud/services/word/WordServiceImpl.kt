package com.guigou.botisgud.services.word

import me.tatarka.inject.annotations.Inject
import java.io.File

@Inject
class WordServiceImpl : WordService {
    // TODO: readLines should not be used for large files
    private val words = File("words").readLines()

    override suspend fun random(): String {
        return words.random()
    }
}
