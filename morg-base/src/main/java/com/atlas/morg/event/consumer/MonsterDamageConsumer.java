package com.atlas.morg.event.consumer;

import com.atlas.csrv.constant.EventConstants;
import com.atlas.csrv.event.MonsterDamageEvent;
import com.atlas.kafka.consumer.SimpleEventHandler;
import com.atlas.morg.MonsterRegistry;
import com.atlas.morg.event.producer.MonsterKilledEventProducer;
import com.atlas.morg.model.DamageSummary;
import com.atlas.morg.model.Monster;

public class MonsterDamageConsumer implements SimpleEventHandler<MonsterDamageEvent> {
   @Override
   public void handle(Long key, MonsterDamageEvent event) {
      MonsterRegistry.getInstance()
            .getMonster(event.uniqueId())
            .filter(Monster::alive)
            .ifPresent(monster -> applyDamage(event.characterId(), event.damage(), monster));

      System.out.printf("Monster %d damaged for %d by %d", event.uniqueId(), event.damage(), event.characterId());
   }

   protected void applyDamage(int characterId, int damage, Monster monster) {
      MonsterRegistry.getInstance()
            .applyDamage(characterId, damage, monster.uniqueId())
            .ifPresent(summary -> {
               // TODO broadcast HP bar update
               if (summary.killed()) {
                  killMonster(monster, summary);
               }
            });
   }

   protected void killMonster(Monster monster, DamageSummary summary) {
      MonsterKilledEventProducer.sendKilled(monster.worldId(), monster.channelId(), monster.mapId(), monster.uniqueId(),
            monster.monsterId(), monster.x(), monster.y(), summary.characterId());
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
      return System.getenv(EventConstants.TOPIC_MONSTER_DAMAGE);
   }
}
