package com.atlas.morg.rest.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.morg.rest.attribute.MonsterAttributes;

import builder.AttributeResultBuilder;

public class MonsterAttributesBuilder extends RecordBuilder<MonsterAttributes, MonsterAttributesBuilder>
      implements AttributeResultBuilder {
   private Integer monsterId;

   @Override
   public MonsterAttributes construct() {
      return new MonsterAttributes(monsterId);
   }

   @Override
   public MonsterAttributesBuilder getThis() {
      return this;
   }

   public MonsterAttributesBuilder setMonsterId(Integer monsterId) {
      this.monsterId = monsterId;
      return getThis();
   }
}
