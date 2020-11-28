package com.guigou.botisgud.models

import org.bson.codecs.pojo.annotations.BsonId
import org.litote.kmongo.Id
import org.litote.kmongo.newId
import java.time.Instant

// TODO: can Snowflake be serialized? if so then convert userId back to Snowflake.
data class Reminder(
    val userId: Long,
    val message: String,
    val link: String,
    val timestamp: Instant
) {
    @BsonId val key: Id<Reminder> = newId()
}
