package com.guigou.botisgud

import com.gitlab.kordlib.core.Kord
import com.guigou.botisgud.commands.Nicknamer
import com.guigou.botisgud.commands.Reminder
import com.guigou.botisgud.commands.Typing
import com.guigou.botisgud.commands.register
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val token = System.getenv("DISCORD_TOKEN")
    val client = Kord(token)

    // https://elizarov.medium.com/coroutine-context-and-scope-c8b255d59055
    client.register(Typing(), this) // TODO: determine how to initialize class within extension function
    client.register(Reminder(), this)
    client.register(Nicknamer(), this)

    client.login()
}
