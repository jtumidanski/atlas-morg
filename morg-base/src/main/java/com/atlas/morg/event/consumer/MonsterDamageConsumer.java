package com.atlas.morg.event.consumer;

import com.atlas.csrv.constant.EventConstants;
import com.atlas.csrv.event.MonsterDamageEvent;
import com.atlas.kafka.consumer.SimpleEventHandler;
import com.atlas.morg.MonsterRegistry;
import com.atlas.morg.event.producer.MonsterKilledEventProducer;
import com.atlas.morg.model.DamageEntry;
import com.atlas.morg.model.DamageSummary;
import com.atlas.morg.model.Monster;
import com.atlas.morg.processor.MonsterProcessor;
import com.atlas.morg.processor.TopicDiscoveryProcessor;

public class MonsterDamageConsumer implements SimpleEventHandler<MonsterDamageEvent> {
   @Override
   public void handle(Long key, MonsterDamageEvent event) {
      MonsterRegistry.getInstance()
            .getMonster(event.uniqueId())
            .filter(Monster::alive)
            .ifPresent(monster -> applyDamage(event.characterId(), event.damage(), monster));

      System.out.printf("Monster %d damaged for %d by %d", event.uniqueId(), event.damage(), event.characterId());
   }

   protected void applyDamage(int characterId, long damage, Monster monster) {
      MonsterRegistry.getInstance()
            .applyDamage(characterId, damage, monster.uniqueId())
            .ifPresent(summary -> {
               if (characterId != summary.monster().controlCharacterId()) {
                  boolean damageLeader = summary.monster().damageLeader()
                        .map(DamageEntry::characterId)
                        .filter(id -> id == characterId)
                        .isPresent();
                  if (damageLeader) {
                     MonsterProcessor.stopControl(summary.monster());
                     MonsterProcessor.startControl(summary.monster(), characterId);
                  }
               }

               // TODO broadcast HP bar update
               if (summary.killed()) {
                  killMonster(monster, summary);
               }
            });
   }

   protected void killMonster(Monster monster, DamageSummary summary) {
      MonsterKilledEventProducer.sendKilled(monster.worldId(), monster.channelId(), monster.mapId(), monster.uniqueId(),
            monster.monsterId(), monster.x(), monster.y(), summary.characterId(), monster.damageSummary());
      MonsterRegistry.getInstance().removeMonster(monster.uniqueId());
   }

   @Override
   public Class<MonsterDamageEvent> getEventClass() {
      return MonsterDamageEvent.class;
   }

   @Override
   public String getConsumerId() {
      return "Monster Registry";
   }

   @Override
   public String getBootstrapServers() {
      return System.getenv("BOOTSTRAP_SERVERS");
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_MONSTER_DAMAGE);
   }
}
