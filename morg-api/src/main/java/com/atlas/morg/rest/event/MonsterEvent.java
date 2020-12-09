package com.atlas.morg.rest.event;

public record MonsterEvent(int worldId, int channelId, int mapId, int uniqueId, MonsterEventType type) {
}
