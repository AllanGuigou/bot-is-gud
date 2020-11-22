package com.guigou.botisgud.commands

import com.gitlab.kordlib.common.entity.Snowflake
import com.gitlab.kordlib.core.Kord
import com.gitlab.kordlib.core.behavior.edit
import com.gitlab.kordlib.core.event.message.MessageCreateEvent
import com.gitlab.kordlib.core.on
import com.guigou.botisgud.services.WordService
import com.guigou.botisgud.services.WordServiceImpl
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.launch
import java.time.Instant
import java.util.*
import kotlin.concurrent.schedule

class Nicknamer(private val wordService: WordService = WordServiceImpl()) : Command {

    private val users: MutableMap<Snowflake, MutableSet<Snowflake>> = mutableMapOf()

    override fun register(client: Kord) {
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

        // TODO: would this be a daemon?
        val time = Date.from(Instant.now().plusSeconds(60))!!
        val period = 5 * 60 * 1000L
        Timer("nicknamer").schedule(time, period) {
            GlobalScope.launch {
                for (entry in users) {
                    for (userSnowflake in entry.value) {
                        val member = client.getUser(userSnowflake)!!.asMember(entry.key)

                        member.edit {
                            nickname = wordService.random()
                        }
                    }
                }
            }
        }
    }
}
