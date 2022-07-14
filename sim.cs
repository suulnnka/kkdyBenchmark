
using System;
using System.Diagnostics;
using System.Threading.Tasks;

namespace Suulnnka {

  class Random_count {
    const int MBIG = Int32.MaxValue;
    const int MSEED = 161803398;

    const int att = 446;
    const int att_ = 714;

    const int att_s = 624;
    const int att_s_ = 999;

    const double prob = 0.2;

    public int seed;
    int mp = 0;

    int inext, inextp;
    int[] SeedArray = new int[56];

    public Random_count(int Seed) {
      this.seed = Seed;

      int ii;
      int mj, mk;

      // Numerical Recipes in C online @ http://www.library.cornell.edu/nr/bookcpdf/c7-1.pdf
      mj = MSEED - Math.Abs(Seed);
      SeedArray[55] = mj;
      mk = 1;
      for (int i = 1; i < 55; i++) {  //  [1, 55] is special (Knuth)
        ii = (21 * i) % 55;
        SeedArray[ii] = mk;
        mk = mj - mk;
        if (mk < 0)
          mk += MBIG;
        mj = SeedArray[ii];
      }
      for (int k = 1; k < 5; k++) {
        for (int i = 1; i < 56; i++) {
          SeedArray[i] -= SeedArray[1 + (i + 30) % 55];
          if (SeedArray[i] < 0)
            SeedArray[i] += MBIG;
        }
      }
      inext = 0;
      inextp = 31;
    }

    public Random_count(Random_count r) {
      seed = r.seed;
      mp = r.mp;
      inext = r.inext;
      inextp = r.inextp;
      for (int i = 0; i < SeedArray.Length; i++) {
        SeedArray[i] = r.SeedArray[i];
      }
    }

    public double NextDouble() {
      int retVal;

      if (++inext >= 56)
        inext = 1;
      if (++inextp >= 56)
        inextp = 1;

      retVal = SeedArray[inext] - SeedArray[inextp];

      if (retVal < 0)
        retVal += MBIG;

      SeedArray[inext] = retVal;

      return retVal * (1.0 / MBIG);
    }

    public int attack(int def, int hp) {
      bool is_skill = mp == 4;
      int d;
      if (is_skill) {
        if (NextDouble() > prob) {
          d = att_s - def;
        } else {
          d = att_s_ - def;
        }

        if (hp > d) {
          if (NextDouble() > prob) {
            d += att_s - def;
          } else {
            d += att_s_ - def;
          }
        }

      } else {
        if (NextDouble() > prob) {
          d = att - def;
        } else {
          d = att_ - def;
        }
      }

      mp++;
      if (mp > 4) {
        mp = 0;
      }
      return d;
    }

  }

  class Kkdy112 {

    static bool debug = false;

    static void Main() {

      var start = DateTime.Now;

      int l = 0;

      var for_times = 2150;
      var for_step = 1000000L;
      if (true) {
        for_times = 0;
        for_step = 1000L;
      }

      //*
      Parallel.For(0, for_times, i => {
        long step = for_step;
        long min = i * step;
        long max = (i + 1) * step;
        if (min < Int32.MinValue)
          min = Int32.MinValue;
        if (max < Int32.MinValue)
          max = Int32.MinValue;
        if (min > Int32.MaxValue)
          min = Int32.MaxValue;
        if (max > Int32.MaxValue)
          max = Int32.MaxValue;

        for (int j = (int)min; j < (int)max; j++) {
          run(j);
        }

        l++;
        print(l+"/"+ for_times);
      });
      //*/

      // run(4050455);

      print((DateTime.Now - start).ToString());
      /*
      print("漏12狗 " + possible[0]);
      print("漏17狗 " + possible[1]);
      print("漏16双 " + possible[2]);
      print("漏18狗 " + possible[3]);
      print("漏19双 " + possible[4]);
      print("漏20狗 " + possible[5]);
      print("漏21双 " + possible[6]);
      print("漏23投 " + possible[7]);

      print("成功 " + possible[99]);

      for (int i = 8; i < 30; i++) {
        if (possible[i] > 0) {
          print(i + " " + possible[i]);
        }
      }
      */
      if (!debug) {
        Console.ReadKey();
        Console.ReadKey();
      }
    }

    static void print(string str) {
      if (debug) {
        Trace.WriteLine(str);
      } else {
        Console.WriteLine(str);
      }
    }

    // static int[] possible = new int[100];

    static bool run(int seed) {
      var r = new Random_count(seed);
      

      // 第一阶段 开局 6 狗
      {
        for (int i = 1; i <= 6; i++) {
          int hp = hp_g;

          while (hp > 0) {
            int d = r.attack(def_g,hp);
            hp = hp - d;
          }
        }
      }

      // 第二阶段 7-10 双 狗 狗 狗
      {
        int hp_7 = hp_s;
        for (int i = 0; i < 4; i++) {
          hp_7 -= r.attack(def_s,hp_7);
          if (hp_7 <= 0)
            break;
        }

        int hp_8 = hp_g;
        for (int i = 0; i < 4; i++) {
          hp_8 -= r.attack(def_g,hp_8);
          if (hp_8 <= 0)
            break;
        }

        if (hp_7 > 0) {
          for (int i = 0; i < 2; i++) {
            hp_7 -= r.attack(def_s,hp_7);
            if (hp_7 <= 0)
              break;
          }
        }

        int hp_9 = hp_g;
        for (int i = 0; i < 4; i++) {
          hp_9 -= r.attack(def_g,hp_9);
          if (hp_9 <= 0)
            break;
        }

        int hp_10 = hp_g;
        for (int i = 0; i < 4; i++) {
          hp_10 -= r.attack(def_g,hp_10);
          if (hp_10 <= 0)
            break;
        }
      }

      // 第三阶段 11 - 13 双 狗 盾
      {

        // 11 双 4
        // 12 狗 3
        // 11 双 1

        // 有7下分配到 11，12两怪
        int hp_11 = hp_s;
        int hp_12 = hp_g;
        int hp_13 = hp_d;

        int count_3 = 0;

        while(true) {
          hp_11 -= r.attack(def_s,hp_11);
          count_3++;
          if (hp_11 <= 0)
            break;
        }

        if (count_3 == 2) {
          hp_13 -= r.attack(def_d,hp_13);
          count_3++;
        }

        while (true) {
          hp_12 -= r.attack(def_g,hp_12);
          count_3++;
          if (hp_12 <= 0)
            break;
        }

        if (count_3 > 7) {
          //possible[0]++;
          return false;
        }

        // 至少能打8下，13怪必死
        for (int i = count_3; i < 15; i++) {
          hp_13 -= r.attack(def_d,hp_13);
          if (hp_13 <= 0)
            break;
        }
      }

      // 第四阶段 14 - 17 投 狗 双 狗
      int count_4 = 0;
      {
        // 14 投 5
        // 15 狗 4
        // 16 双 接上面的 8 
        // 17 狗 到12
        // 16 双 到17
        int hp_14 = hp_t;
        for (int i = 0; i < 5; i++) {
          hp_14 -= r.attack(def_t,hp_14);
          count_4++;
          if (hp_14 <= 0)
            break;
        }

        int hp_15 = hp_g;
        for (int i = 0; i < 4; i++) {
          hp_15 -= r.attack(def_g,hp_15);
          count_4++;
          if (hp_15 <= 0)
            break;
        }

        // 与上两个怪均分8下攻击
        int hp_16 = hp_s;
        for (int i = count_4; i < 8; i++) {
          hp_16 -= r.attack(def_s,hp_16);
          count_4++;
          if (hp_16 <= 0)
            break;
        }

        // 如果上三个怪8下死了，那么就不会漏 17 狗
        
        int hp_17 = hp_g;
        for (int i = count_4; i < 12; i++) {
          hp_17 -= r.attack(def_g,hp_17);
          count_4++;
          if (hp_17 <= 0)
            break;
        }

        if (hp_17 > 0) {
          //possible[1]++;
          return false;
        }

        // 如果16没死，将获得17次攻击以前（含17）的全部剩余攻击
        if (hp_16 > 0) {
          for (int i = count_4; i < 17; i++) {
            hp_16 -= r.attack(def_s,hp_16);
            count_4++;
            if (hp_16 <= 0)
              break;
          }
        }
        if (hp_16 > 0) {
          //possible[2]++;
          return false;
        }
        
      }

      // 第五阶段
      // 如果四阶段17下解决，状况1，攻击次数正常
      // 如果四阶段16下解决，状况1，攻击次数+1
      // 如果四阶段15下解决，状况1，攻击次数+2
      // 如果四阶段14下及以内解决，状况2
      // 如果四阶段13下及以内解决，状况2

      if (count_4 >= 15) {
        // 状况 1 
        // 18 狗 2 3 4
        // 19 双 5 5 5
        // 20 狗 0 0 0
        // 21 双 5 5 5

        return false;

        /*

        int count_5 = count_4 - 15;

        int hp_18 = hp_g;
        for (int i = count_5; i < 4; i++) {
          hp_18 -= r.attack(def_g);
          count_5++;
          if (hp_18 <= 0)
            break;
        }
        if (hp_18 > 0) {
          //possible[3]++;
          return false;
        }

        int hp_19 = hp_s;
        for (int i = count_5; i < 9; i++) {
          hp_19 -= r.attack(def_s);
          count_5++;
          if (hp_19 <= 0)
            break;
        }
        if (hp_19 > 0) {
          //possible[4]++;
          return false;
        }

        int hp_20 = hp_g;
        for (int i = count_5; i < 9; i++) {
          hp_20 -= r.attack(def_g);
          count_5++;
          if (hp_20 <= 0)
            break;
        }
        if (hp_20 > 0) {
          //possible[5]++;
          return false;
        }

        int hp_21 = hp_s;
        for (int i = count_5; i < 14; i++) {
          hp_21 -= r.attack(def_s);
          count_5++;
          if (hp_21 <= 0)
            break;
        }
        if (hp_21 > 0) {
          //possible[6]++;
          return false;
        }

        // TODO 
        return false;

        */

      } else {

        // 状况 2
        // 18 狗 4
        // 19 双 5
        // 20 狗 1
        // 21 双 3
        int count_5 = 0;
        //阶段5-1
        {
          int hp_18 = hp_g;
          for (int i = count_5; i < 4; i++) {
            hp_18 -= r.attack(def_g,hp_18);
            count_5++;
            if (hp_18 <= 0)
              break;
          }
          if (hp_18 > 0) {
            //possible[3]++;
            return false;
          }

          int hp_19 = hp_s;
          for (int i = count_5; i < 9; i++) {
            hp_19 -= r.attack(def_s,hp_19);
            count_5++;
            if (hp_19 <= 0)
              break;
          }
          if (hp_19 > 0) {
            //possible[4]++;
            return false;
          }

          if (count_5 <= 5) {
            // 另外一种轴 ? 
            // TODO
            return false;
          }

          int hp_20 = hp_g;
          for (int i = count_5; i < 10; i++) {
            hp_20 -= r.attack(def_g,hp_20);
            count_5++;
            if (hp_20 <= 0)
              break;
          }
          if (hp_20 > 0) {
            //possible[5]++;
            return false;
          }

        }

        // 阶段5-2

        // 12下会乱轴

        // 21 打3下
        // 22 打3下
        // 23 打2下

        var r_pos_2 = new Random_count(r);
        var r_pos_3 = new Random_count(r);
        var r_pos_4 = new Random_count(r);


        part_5_2_pos_2(count_5, r_pos_2);
        part_5_2_pos_3(count_5, r_pos_3);
        part_5_2_pos_4(count_5, r_pos_4);

        return false;
      }

      // possible[99]++;
      // return false;

    }

    static bool part_5_2_pos_4(int count_5, Random_count r) {

      // 22 双 及 23 投 进入攻击范围时启动装置

      int hp_21 = hp_s;
      for (int i = count_5; i < 13; i++) {
        hp_21 -= r.attack(def_s, hp_21);
        count_5++;
        if (hp_21 <= 0)
          break;
      }

      if (count_5 <= 12) {
        // print("???????? TODO");
        // 5-2-2 轴 TODO
        return false;
      }
      int hp_22 = hp_s;
      for (int i = count_5; i < 14; i++) {
        hp_22 -= r.attack(def_s, hp_22);
        count_5++;
        if (hp_22 <= 0)
          break;
      }

      hp_22 = hp_22 - 1000;

      if (hp_22 > 0) {
        for (int i = count_5; i < 16; i++) {
          hp_22 -= r.attack(def_s, hp_22);
          count_5++;
          if (hp_22 <= 0)
            break;
        }
      }

      if (hp_22 > 0) {
        return false;
      }

      int hp_23 = hp_t - 1000;
      for (int i = count_5; i < 17; i++) {
        hp_23 -= r.attack(def_t, hp_23);
        count_5++;
        if (hp_23 <= 0)
          break;
      }

      if (hp_23 > 0) {
        //possible[7]++;
        return false;
      }

      int hp_24 = hp_s;
      for (int i = count_5; i < 18; i++) {
        hp_24 -= r.attack(def_s, hp_24);
        count_5++;
        if (hp_24 <= 0)
          break;
      }

      int hp_25 = hp_w;
      for (int i = count_5; i < 25; i++) {
        hp_25 -= r.attack(def_w, hp_25);
        count_5++;
      }

      for (int i = count_5; i < 27; i++) {
        hp_24 -= r.attack(def_s, hp_24);
        count_5++;
        if (hp_24 <= 0)
          break;
      }

      for (int i = count_5; i < 27; i++) {
        hp_25 -= r.attack(def_w, hp_25);
        count_5++;
      }

      int count_last = 0;

      if (hp_24 > 0) {
        for (int i = count_last; i < 2; i++) {
          hp_24 -= r.attack(def_s, hp_24);
          count_last++;
          if (hp_24 <= 0)
            break;
        }
      }

      if (hp_24 > 0) {
        return false;
      }

      int hp_26 = hp_s;
      for (int i = count_last; i < 5; i++) {
        hp_26 -= r.attack(def_s, hp_26);
        count_last++;
        if (hp_26 <= 0)
          break;
      }

      for (int i = count_last; i < 8; i++) {
        hp_25 -= r.attack(def_w, hp_25);
        count_last++;
        if (hp_25 <= 0)
          break;
      }

      int hp_27 = hp_g;
      for (int i = count_last; i < 11; i++) {
        hp_27 -= r.attack(def_g, hp_27);
        count_last++;
        if (hp_27 <= 0)
          break;
      }

      if (hp_27 > 0) {
        //possible[8]++;
        return false;
      }

      for (int i = count_last; i < 11; i++) {
        hp_25 -= r.attack(def_w, hp_25);
        count_last++;
        if (hp_25 <= 0)
          break;
      }

      int hp_29 = hp_l;
      for (int i = count_last; i < 15; i++) {
        hp_29 -= r.attack(def_l, hp_29);
        count_last++;
        if (hp_29 <= 0)
          break;
      }

      if (hp_25 > 0) {
        for (int i = count_last; i < 17; i++) {
          hp_25 -= r.attack(def_w, hp_25);
          count_last++;
          if (hp_25 <= 0)
            break;
        }
      }

      if (hp_25 <= 0) {
        // print("w done. seed:" + r.seed);
      }

      if (hp_25 > 0) {
        return false;
      }

      for (int i = count_last; i < 20; i++) {
        hp_29 -= r.attack(def_l, hp_29);
        count_last++;
        if (hp_29 <= 0)
          break;
      }

      // print("流 hp:" + hp_29 + ", seed:" + r.seed);

      if (hp_29 > 0) {
        return false;
      }

      int hp_28 = hp_d;
      for (int i = count_last; i < 23; i++) {
        hp_28 -= r.attack(def_d, hp_28);
        count_last++;
        if (hp_28 <= 0)
          break;
      }

      print("盾 hp:" + hp_28 + ", seed:" + r.seed);

      if (hp_28 > 0) {
        return false;
      }

      print("all done. seed:" + r.seed);

      return true;
    }

    static bool part_5_2_pos_3(int count_5, Random_count r) {

      int hp_21 = hp_s;
      for (int i = count_5; i < 13; i++) {
        hp_21 -= r.attack(def_s, hp_21);
        count_5++;
        if (hp_21 <= 0)
          break;
      }

      if (count_5 <= 12) {
        // print("???????? TODO");
        // 5-2-2 轴 TODO
        return false;
      }
      int hp_22 = hp_s;
      for (int i = count_5; i < 16; i++) {
        hp_22 -= r.attack(def_s, hp_22);
        count_5++;
        if (hp_22 <= 0)
          break;
      }

      if (hp_22 > 0) {
        return false;
      }

      if (hp_21 > 0) {
        for (int i = count_5; i < 16; i++) {
          hp_21 -= r.attack(def_s, hp_21);
          count_5++;
          if (hp_21 <= 0)
            break;
        }
      }

      int hp_23 = hp_t;
      for (int i = count_5; i < 17; i++) {
        hp_23 -= r.attack(def_t,hp_23);
        count_5++;
        if (hp_23 <= 0)
          break;
      }

      hp_23 = hp_23 - 1000;
      if (hp_23 > 0) {
        //for (int i = count_5; i < 18; i++) {
          hp_23 -= r.attack(def_t,hp_23);
          count_5++;
          //if (hp_23 <= 0)
            //break;
        //}
      }

      if (hp_23 > 0) {
        //possible[7]++;
        return false;
      }

      int hp_24 = hp_s - 1000;
      for (int i = count_5; i < 18; i++) {
        hp_24 -= r.attack(def_s,hp_24);
        count_5++;
        if (hp_24 <= 0)
          break;
      }

      int hp_25 = hp_w;
      for (int i = count_5; i < 27; i++) {
        hp_25 -= r.attack(def_w,hp_25);
        count_5++;
      }



      int count_last = 0;

      int hp_26 = hp_s;
      for (int i = count_last; i < 5; i++) {
        hp_26 -= r.attack(def_s,hp_26);
        count_last++;
        if (hp_26 <= 0)
          break;
      }

      for (int i = count_last; i < 5; i++) {
        hp_25 -= r.attack(def_w,hp_25);
        count_last++;
        if (hp_25 <= 0)
          break;
      }

      if (hp_24 > 0) {
        for (int i = count_last; i < 8; i++) {
          hp_24 -= r.attack(def_s,hp_24);
          count_last++;
          if (hp_24 <= 0)
            break;
        }
      }

      for (int i = count_last; i < 8; i++) {
        hp_25 -= r.attack(def_w,hp_25);
        count_last++;
        if (hp_25 <= 0)
          break;
      }

      int hp_27 = hp_g;
      for (int i = count_last; i < 11; i++) {
        hp_27 -= r.attack(def_g,hp_27);
        count_last++;
        if (hp_27 <= 0)
          break;
      }

      if (hp_27 > 0) {
        //possible[8]++;
        return false;
      }

      for (int i = count_last; i < 11; i++) {
        hp_25 -= r.attack(def_w,hp_25);
        count_last++;
        if (hp_25 <= 0)
          break;
      }

      int hp_29 = hp_l;
      for (int i = count_last; i < 15; i++) {
        hp_29 -= r.attack(def_l,hp_29);
        count_last++;
        if (hp_29 <= 0)
          break;
      }

      for (int i = count_last; i < 17; i++) {
        hp_25 -= r.attack(def_w,hp_25);
        count_last++;
        if (hp_25 <= 0)
          break;
      }

      if (hp_25 <= 0) {
        // print("w done. seed:" + r.seed);
      }

      if (hp_25 > 0) {
        return false;
      }

      for (int i = count_last; i < 20; i++) {
        hp_29 -= r.attack(def_l,hp_29);
        count_last++;
        if (hp_29 <= 0)
          break;
      }

      // print("流 hp:"+ hp_29 +", seed:" + r.seed);

      if (hp_29 > 0) {
        return false;
      }

      int hp_28 = hp_d;
      for (int i = count_last; i < 23; i++) {
        hp_28 -= r.attack(def_d,hp_28);
        count_last++;
        if (hp_28 <= 0)
          break;
      }

      print("盾 hp:" + hp_28 + ", seed:" + r.seed);

      if (hp_28 > 0) {
        return false;
      }

      print("all done. seed:" + r.seed);

      return true;
    }
    
    static bool part_5_2_pos_2(int count_5, Random_count r) {

      int hp_21 = hp_s;
      for (int i = count_5; i < 13; i++) {
        hp_21 -= r.attack(def_s,hp_21);
        count_5++;
        if (hp_21 <= 0)
          break;
      }

      if (hp_21 <= 0) {
        // TODO 又乱轴了
        // return false;
      }

      if (count_5 <= 12) {
        // print("???????? TODO");
        // 5-2-2 轴 TODO
        return false;
      }
      int hp_22 = hp_s;
      for (int i = count_5; i < 16; i++) {
        hp_22 -= r.attack(def_s, hp_22);
        count_5++;
        if (hp_22 <= 0)
          break;
      }

      if (hp_22 <= 0) {
        // 又乱轴了!!!
        // TODO
        // return false;
      }

      int hp_23 = hp_t;
      for (int i = count_5; i < 16; i++) {
        hp_23 -= r.attack(def_t, hp_23);
        count_5++;
        if (hp_23 <= 0)
          break;
      }

      hp_23 = hp_23 - 1000;
      for (int i = count_5; i < 18; i++) {
        hp_23 -= r.attack(def_t, hp_23);
        count_5++;
        if (hp_23 <= 0)
          break;
      }

      if (hp_23 > 0) {
        //possible[7]++;
        return false;
      }

      //possible[count_5]++;

      int hp_24 = hp_s - 1000;
      for (int i = count_5; i < 18; i++) {
        hp_24 -= r.attack(def_s, hp_24);
        count_5++;
        if (hp_24 <= 0)
          break;
      }

      int hp_25 = hp_w;
      for (int i = count_5; i < 27; i++) {
        hp_25 -= r.attack(def_w, hp_25);
        count_5++;
      }


      // 5-3阶段
      // 26 双 5
      // 24 双 4
      // 27 狗 2
      // 29 流 4 
      // 25  W 2
      // 29 流 3
      // 28 盾 3
      // 30 投 7


      // 流 W 
      // 流 流
      // 流 流
      // 流 流
      // W  W
      // W  流
      // 流 流
      // 流 流
      // 流 流
      // 盾 盾
      // 盾 盾
      // 盾 盾

      // 5-3阶段
      // 26 双 4 3
      // 25  W 1 1
      // 24 双 3 3
      // 27 狗 2 1
      // 25  W 1 3
      // 29 流 3 3
      // 25  W 2 1
      // 29 流 3 4
      // 28 盾 3 3

      var r_bak = new Random_count(r);

      {
        //* 标准轴
        int count_last = 0;

        int hp_26 = hp_s;
        for (int i = count_last; i < 5; i++) {
          hp_26 -= r.attack(def_s, hp_26);
          count_last++;
          if (hp_26 <= 0)
            break;
        }

        for (int i = count_last; i < 5; i++) {
          hp_25 -= r.attack(def_w, hp_25);
          count_last++;
          if (hp_25 <= 0)
            break;
        }

        if (hp_24 > 0) {
          for (int i = count_last; i < 8; i++) {
            hp_24 -= r.attack(def_s, hp_24);
            count_last++;
            if (hp_24 <= 0)
              break;
          }
        } else {
          hp_25 -= r.attack(def_w, hp_25);
          count_last++;
        }


        int hp_27 = hp_g;
        for (int i = count_last; i < 11; i++) {
          hp_27 -= r.attack(def_g, hp_27);
          count_last++;
          if (hp_27 <= 0)
            break;
        }

        if (hp_27 > 0) {
          //possible[8]++;
          return false;
        }

        for (int i = count_last; i < 11; i++) {
          hp_25 -= r.attack(def_w, hp_25);
          count_last++;
          if (hp_25 <= 0)
            break;
        }

        int hp_29 = hp_l;
        for (int i = count_last; i < 14; i++) {
          hp_29 -= r.attack(def_l, hp_29);
          count_last++;
          if (hp_29 <= 0)
            break;
        }

        for (int i = count_last; i < 16; i++) {
          hp_25 -= r.attack(def_w, hp_25);
          count_last++;
          if (hp_25 <= 0)
            break;
        }

        if (hp_25 <= 0) {
          // print("w done. seed:" + r.seed);
        }

        if (hp_25 > 0) {
          // possible[11]++;
          return false;
        }

        for (int i = count_last; i < 19; i++) {
          hp_29 -= r.attack(def_l, hp_29);
          count_last++;
          if (hp_29 <= 0)
            break;
        }

        // print("流hp: " + hp_29 + ", seed:"+r.seed);
        if (hp_29 > 0) {
          // possible[9]++;
          return false;
        }

        int hp_28 = hp_d;
        for (int i = count_last; i < 22; i++) {
          hp_28 -= r.attack(def_d, hp_28);
          count_last++;
          if (hp_28 <= 0)
            break;
        }

        print("盾hp: " + hp_28 + ", seed:" + r.seed);

        if (hp_28 > 0) {
          // possible[10]++;
          return false;
        }

        print("all done. seed:" + r.seed);

        // return true;

        //*/
      }

      r = r_bak;
      //*
      {
        int count_last = 0;

        int hp_26 = hp_s;
        for (int i = count_last; i < 5; i++) {
          hp_26 -= r.attack(def_s, hp_26);
          count_last++;
          if (hp_26 <= 0)
            break;
        }

        for (int i = count_last; i < 4; i++) {
          hp_25 -= r.attack(def_w,hp_25);
          count_last++;
          if (hp_25 <= 0)
            break;
        }

        if (hp_24 > 0) {
          for (int i = count_last; i < 8; i++) {
            hp_24 -= r.attack(def_s, hp_24);
            count_last++;
            if (hp_24 <= 0)
              break;
          }
        } else {
          hp_25 -= r.attack(def_w, hp_25);
          count_last++;
        }


        int hp_27 = hp_g;
        for (int i = count_last; i < 11; i++) {
          hp_27 -= r.attack(def_g, hp_27);
          count_last++;
          if (hp_27 <= 0)
            break;
        }

        if (hp_27 > 0) {
          //possible[8]++;
          return false;
        }

        for (int i = count_last; i < 11; i++) {
          hp_25 -= r.attack(def_w, hp_25);
          count_last++;
          if (hp_25 <= 0)
            break;
        }

        int hp_29 = hp_l;
        for (int i = count_last; i < 14; i++) {
          hp_29 -= r.attack(def_l, hp_29);
          count_last++;
          if (hp_29 <= 0)
            break;
        }

        for (int i = count_last; i < 15; i++) {
          hp_25 -= r.attack(def_w, hp_25);
          count_last++;
          if (hp_25 <= 0)
            break;
        }

        if (hp_25 <= 0) {
          // print("w done. seed:" + r.seed);
        }

        if (hp_25 > 0) {
          // possible[11]++;
          return false;
        }

        for (int i = count_last; i < 19; i++) {
          hp_29 -= r.attack(def_l, hp_29);
          count_last++;
          if (hp_29 <= 0)
            break;
        }

        // print("流hp: " + hp_29 + ", seed:" + r.seed);
        if (hp_29 > 0) {
          // possible[9]++;
          return false;
        }

        int hp_28 = hp_d;
        for (int i = count_last; i < 22; i++) {
          hp_28 -= r.attack(def_d, hp_28);
          count_last++;
          if (hp_28 <= 0)
            break;
        }

        print("盾hp: " + hp_28 + ", seed:" + r.seed);

        if (hp_28 > 0) {
          // possible[10]++;
          return false;
        }

        print("all done. seed:" + r.seed);

        return true;
      }
      //*/
    }
  }
}
