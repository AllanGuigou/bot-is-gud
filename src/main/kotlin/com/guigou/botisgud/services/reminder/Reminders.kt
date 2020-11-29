package com.guigou.botisgud.services.reminder

import org.jetbrains.exposed.dao.id.IntIdTable
import org.jetbrains.exposed.sql.`java-time`.datetime

object Reminders : IntIdTable() {
    val userId = long("userId")
    val message = varchar("message", 2000)
    val link = varchar("link", 100) // TODO: what is an appropriate length
    val timestamp = datetime("timestamp")
}
