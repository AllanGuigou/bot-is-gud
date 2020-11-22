package com.guigou.botisgud.commands

import com.gitlab.kordlib.common.entity.Snowflake
import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.behavior.edit
import com.gitlab.kordlib.core.event.message.MessageCreateEvent
import com.gitlab.kordlib.core.on
import com.guigou.botisgud.services.WordService
import com.guigou.botisgud.services.WordServiceImpl
import kotlinx.coroutines.*
import java.time.Instant.now

class Nicknamer(private val wordService: WordService = WordServiceImpl()) : Command {

    private val users: MutableMap<Snowflake, MutableSet<Snowflake>> = mutableMapOf()

    override suspend fun register(client: Kord, scope: CoroutineScope) {
        client.on<MessageCreateEvent> {
            if (message.content != "!nicknamer") {
                return@on
            }

            val guild = message.getGuild().id
            val member = message.author!!.asUser().id

            if (users.containsKey(guild)) {
                users[guild]!!.add(member)
            } else {
                users[guild] = mutableSetOf(member)
            }
        }

        backgroundWork(client, scope)
    }

    private fun backgroundWork(client: Kord, scope: CoroutineScope) {
        var time = now().plusSeconds(60)!!
        val period = 5 * 60 * 1000L

        scope.launch(Dispatchers.Default) {
            while (true) {
                delay(time.toEpochMilli() - now().toEpochMilli())
                for (entry in users) {
                    for (userSnowflake in entry.value) {
                        val member = client.getUser(userSnowflake)!!.asMember(entry.key)

                        try {
                            member.edit {
                                nickname = wordService.random()
                            }
                        } catch(ex: Exception) { }
                    }
                }
                time = now().plusMillis(period)
            }
        }
    }

}
