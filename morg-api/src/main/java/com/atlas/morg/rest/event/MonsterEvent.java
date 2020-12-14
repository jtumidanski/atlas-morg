package com.atlas.morg.rest.event;

public record MonsterEvent(int worldId, int channelId, int mapId, int uniqueId, Integer actorId, MonsterEventType type) {
}
