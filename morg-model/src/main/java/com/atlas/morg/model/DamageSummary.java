package com.atlas.morg.model;

public record DamageSummary(int characterId, int uniqueId, long visibleDamage, long actualDamage, boolean killed) {
}
