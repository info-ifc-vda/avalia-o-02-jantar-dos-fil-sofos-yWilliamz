## Deadlock

O deadlock foi resolvido implementando uma estratégia onde os filósofos adquirem os garfos um de cada vez. Se um filósofo não conseguir pegar os dois garfos necessários dentro de um tempo limite, ele libera qualquer garfo que tenha pegado e espera antes de tentar novamente. Isso evita que todos os filósofos fiquem travados simultaneamente esperando por recursos, pois garante que sempre haja movimento na aquisição e liberação dos garfos, prevenindo o bloqueio circular.

## Starvation

O problema de starvation foi evitado ao monitorar o tempo que cada filósofo passa sem comer. Se um filósofo ficar inativo por muito tempo, ele recebe prioridade para acessar os garfos nas próximas tentativas. Esse mecanismo garante que nenhum filósofo permaneça faminto indefinidamente, assegurando que todos tenham a oportunidade de se alimentar de forma justa ao longo do tempo.
