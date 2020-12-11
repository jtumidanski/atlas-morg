package com.atlas.morg.model;

public record DamageSummary(int characterId, int uniqueId, int visibleDamage, int actualDamage, boolean killed) {
}
