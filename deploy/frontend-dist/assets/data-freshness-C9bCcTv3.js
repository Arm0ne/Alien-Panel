import{Bi as e,Fi as t,Gi as n,Hn as r,Kr as i,Li as a,Mi as o,Ni as s,Nr as c,On as l,Pi as u,Ri as d,Tn as f,Vn as p,Xi as m,Yi as h,Yn as g,Zr as _,ai as v,cn as y,da as b,ga as x,ii as S,ji as C,kt as w,la as T,ni as E,nr as D,oi as O,pn as k,qi as A,ra as j,ri as M,si as N,ur as P,zi as F}from"./router-CTymuOix.js";import{d as I}from"./index-hZFH7S1t.js";function L(e){let{textColor2:t,primaryColorHover:n,primaryColorPressed:r,primaryColor:a,infoColor:o,successColor:s,warningColor:c,errorColor:l,baseColor:u,borderColor:d,opacityDisabled:f,tagColor:p,closeIconColor:m,closeIconColorHover:h,closeIconColorPressed:g,borderRadiusSmall:_,fontSizeMini:v,fontSizeTiny:y,fontSizeSmall:b,fontSizeMedium:x,heightMini:S,heightTiny:C,heightSmall:w,heightMedium:T,closeColorHover:E,closeColorPressed:D,buttonColor2Hover:O,buttonColor2Pressed:k,fontWeightStrong:A}=e;return Object.assign(Object.assign({},I),{closeBorderRadius:_,heightTiny:S,heightSmall:C,heightMedium:w,heightLarge:T,borderRadius:_,opacityDisabled:f,fontSizeTiny:v,fontSizeSmall:y,fontSizeMedium:b,fontSizeLarge:x,fontWeightStrong:A,textColorCheckable:t,textColorHoverCheckable:t,textColorPressedCheckable:t,textColorChecked:u,colorCheckable:`#0000`,colorHoverCheckable:O,colorPressedCheckable:k,colorChecked:a,colorCheckedHover:n,colorCheckedPressed:r,border:`1px solid ${d}`,textColor:t,color:p,colorBordered:`rgb(250, 250, 252)`,closeIconColor:m,closeIconColorHover:h,closeIconColorPressed:g,closeColorHover:E,closeColorPressed:D,borderPrimary:`1px solid ${i(a,{alpha:.3})}`,textColorPrimary:a,colorPrimary:i(a,{alpha:.12}),colorBorderedPrimary:i(a,{alpha:.1}),closeIconColorPrimary:a,closeIconColorHoverPrimary:a,closeIconColorPressedPrimary:a,closeColorHoverPrimary:i(a,{alpha:.12}),closeColorPressedPrimary:i(a,{alpha:.18}),borderInfo:`1px solid ${i(o,{alpha:.3})}`,textColorInfo:o,colorInfo:i(o,{alpha:.12}),colorBorderedInfo:i(o,{alpha:.1}),closeIconColorInfo:o,closeIconColorHoverInfo:o,closeIconColorPressedInfo:o,closeColorHoverInfo:i(o,{alpha:.12}),closeColorPressedInfo:i(o,{alpha:.18}),borderSuccess:`1px solid ${i(s,{alpha:.3})}`,textColorSuccess:s,colorSuccess:i(s,{alpha:.12}),colorBorderedSuccess:i(s,{alpha:.1}),closeIconColorSuccess:s,closeIconColorHoverSuccess:s,closeIconColorPressedSuccess:s,closeColorHoverSuccess:i(s,{alpha:.12}),closeColorPressedSuccess:i(s,{alpha:.18}),borderWarning:`1px solid ${i(c,{alpha:.35})}`,textColorWarning:c,colorWarning:i(c,{alpha:.15}),colorBorderedWarning:i(c,{alpha:.12}),closeIconColorWarning:c,closeIconColorHoverWarning:c,closeIconColorPressedWarning:c,closeColorHoverWarning:i(c,{alpha:.12}),closeColorPressedWarning:i(c,{alpha:.18}),borderError:`1px solid ${i(l,{alpha:.23})}`,textColorError:l,colorError:i(l,{alpha:.1}),colorBorderedError:i(l,{alpha:.08}),closeIconColorError:l,closeIconColorHoverError:l,closeIconColorPressedError:l,closeColorHoverError:i(l,{alpha:.12}),closeColorPressedError:i(l,{alpha:.18})})}var R={name:`Tag`,common:y,self:L},z={color:Object,type:{type:String,default:`default`},round:Boolean,size:String,closable:Boolean,disabled:{type:Boolean,default:void 0}},B=M(`tag`,`
 --n-close-margin: var(--n-close-margin-top) var(--n-close-margin-right) var(--n-close-margin-bottom) var(--n-close-margin-left);
 white-space: nowrap;
 position: relative;
 box-sizing: border-box;
 cursor: default;
 display: inline-flex;
 align-items: center;
 flex-wrap: nowrap;
 padding: var(--n-padding);
 border-radius: var(--n-border-radius);
 color: var(--n-text-color);
 background-color: var(--n-color);
 transition: 
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 line-height: 1;
 height: var(--n-height);
 font-size: var(--n-font-size);
`,[v(`strong`,`
 font-weight: var(--n-font-weight-strong);
 `),S(`border`,`
 pointer-events: none;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 border-radius: inherit;
 border: var(--n-border);
 transition: border-color .3s var(--n-bezier);
 `),S(`icon`,`
 display: flex;
 margin: 0 4px 0 0;
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 font-size: var(--n-avatar-size-override);
 `),S(`avatar`,`
 display: flex;
 margin: 0 6px 0 0;
 `),S(`close`,`
 margin: var(--n-close-margin);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),v(`round`,`
 padding: 0 calc(var(--n-height) / 3);
 border-radius: calc(var(--n-height) / 2);
 `,[S(`icon`,`
 margin: 0 4px 0 calc((var(--n-height) - 8px) / -2);
 `),S(`avatar`,`
 margin: 0 6px 0 calc((var(--n-height) - 8px) / -2);
 `),v(`closable`,`
 padding: 0 calc(var(--n-height) / 4) 0 calc(var(--n-height) / 3);
 `)]),v(`icon, avatar`,[v(`round`,`
 padding: 0 calc(var(--n-height) / 3) 0 calc(var(--n-height) / 2);
 `)]),v(`disabled`,`
 cursor: not-allowed !important;
 opacity: var(--n-opacity-disabled);
 `),v(`checkable`,`
 cursor: pointer;
 box-shadow: none;
 color: var(--n-text-color-checkable);
 background-color: var(--n-color-checkable);
 `,[O(`disabled`,[E(`&:hover`,`background-color: var(--n-color-hover-checkable);`,[O(`checked`,`color: var(--n-text-color-hover-checkable);`)]),E(`&:active`,`background-color: var(--n-color-pressed-checkable);`,[O(`checked`,`color: var(--n-text-color-pressed-checkable);`)])]),v(`checked`,`
 color: var(--n-text-color-checked);
 background-color: var(--n-color-checked);
 `,[O(`disabled`,[E(`&:hover`,`background-color: var(--n-color-checked-hover);`),E(`&:active`,`background-color: var(--n-color-checked-pressed);`)])])])]),V=Object.assign(Object.assign(Object.assign({},f.props),z),{bordered:{type:Boolean,default:void 0},checked:Boolean,checkable:Boolean,strong:Boolean,triggerClickOnClose:Boolean,onClose:[Array,Function],onMouseenter:Function,onMouseleave:Function,"onUpdate:checked":Function,onUpdateChecked:Function,internalCloseFocusable:{type:Boolean,default:!0},internalCloseIsButtonTag:{type:Boolean,default:!0},onCheckedChange:Function}),H=c(`n-tag`),U=F({name:`Tag`,props:V,slots:Object,setup(e){let t=T(null),{mergedBorderedRef:n,mergedClsPrefixRef:i,inlineThemeDisabled:a,mergedRtlRef:o,mergedComponentPropsRef:s}=r(e),c=C(()=>e.size||s?.value?.Tag?.size||`medium`),u=f(`Tag`,`-tag`,B,R,e,i);m(H,{roundRef:b(e,`round`)});function d(){if(!e.disabled&&e.checkable){let{checked:t,onCheckedChange:n,onUpdateChecked:r,"onUpdate:checked":i}=e;r&&r(!t),i&&i(!t),n&&n(!t)}}function h(t){if(e.triggerClickOnClose||t.stopPropagation(),!e.disabled){let{onClose:n}=e;n&&D(n,t)}}let g={setTextContent(e){let{value:n}=t;n&&(n.textContent=e)}},v=l(`Tag`,o,i),y=C(()=>{let{type:t,color:{color:r,textColor:i}={}}=e,a=c.value,{common:{cubicBezierEaseInOut:o},self:{padding:s,closeMargin:l,borderRadius:d,opacityDisabled:f,textColorCheckable:p,textColorHoverCheckable:m,textColorPressedCheckable:h,textColorChecked:g,colorCheckable:v,colorHoverCheckable:y,colorPressedCheckable:b,colorChecked:x,colorCheckedHover:S,colorCheckedPressed:C,closeBorderRadius:w,fontWeightStrong:T,[N(`colorBordered`,t)]:E,[N(`closeSize`,a)]:D,[N(`closeIconSize`,a)]:O,[N(`fontSize`,a)]:k,[N(`height`,a)]:A,[N(`color`,t)]:j,[N(`textColor`,t)]:M,[N(`border`,t)]:P,[N(`closeIconColor`,t)]:F,[N(`closeIconColorHover`,t)]:I,[N(`closeIconColorPressed`,t)]:L,[N(`closeColorHover`,t)]:R,[N(`closeColorPressed`,t)]:z}}=u.value,B=_(l);return{"--n-font-weight-strong":T,"--n-avatar-size-override":`calc(${A} - 8px)`,"--n-bezier":o,"--n-border-radius":d,"--n-border":P,"--n-close-icon-size":O,"--n-close-color-pressed":z,"--n-close-color-hover":R,"--n-close-border-radius":w,"--n-close-icon-color":F,"--n-close-icon-color-hover":I,"--n-close-icon-color-pressed":L,"--n-close-icon-color-disabled":F,"--n-close-margin-top":B.top,"--n-close-margin-right":B.right,"--n-close-margin-bottom":B.bottom,"--n-close-margin-left":B.left,"--n-close-size":D,"--n-color":r||(n.value?E:j),"--n-color-checkable":v,"--n-color-checked":x,"--n-color-checked-hover":S,"--n-color-checked-pressed":C,"--n-color-hover-checkable":y,"--n-color-pressed-checkable":b,"--n-font-size":k,"--n-height":A,"--n-opacity-disabled":f,"--n-padding":s,"--n-text-color":i||M,"--n-text-color-checkable":p,"--n-text-color-checked":g,"--n-text-color-hover-checkable":m,"--n-text-color-pressed-checkable":h}}),x=a?p(`tag`,C(()=>{let t=``,{type:r,color:{color:i,textColor:a}={}}=e;return t+=r[0],t+=c.value[0],i&&(t+=`a${P(i)}`),a&&(t+=`b${P(a)}`),n.value&&(t+=`c`),t}),y,e):void 0;return Object.assign(Object.assign({},g),{rtlEnabled:v,mergedClsPrefix:i,contentRef:t,mergedBordered:n,handleClick:d,handleCloseClick:h,cssVars:a?void 0:y,themeClass:x?.themeClass,onRender:x?.onRender})},render(){var t;let{mergedClsPrefix:n,rtlEnabled:r,closable:i,color:{borderColor:a}={},round:o,onRender:s,$slots:c}=this;s?.();let l=g(c.avatar,t=>t&&e(`div`,{class:`${n}-tag__avatar`},t)),u=g(c.icon,t=>t&&e(`div`,{class:`${n}-tag__icon`},t));return e(`div`,{class:[`${n}-tag`,this.themeClass,{[`${n}-tag--rtl`]:r,[`${n}-tag--strong`]:this.strong,[`${n}-tag--disabled`]:this.disabled,[`${n}-tag--checkable`]:this.checkable,[`${n}-tag--checked`]:this.checkable&&this.checked,[`${n}-tag--round`]:o,[`${n}-tag--avatar`]:l,[`${n}-tag--icon`]:u,[`${n}-tag--closable`]:i}],style:this.cssVars,onClick:this.handleClick,onMouseenter:this.onMouseenter,onMouseleave:this.onMouseleave},u||l,e(`span`,{class:`${n}-tag__content`,ref:`contentRef`},(t=this.$slots).default?.call(t)),!this.checkable&&i?e(k,{clsPrefix:n,class:`${n}-tag__close`,disabled:this.disabled,onClick:this.handleCloseClick,focusable:this.internalCloseFocusable,round:o,isButtonTag:this.internalCloseIsButtonTag,absolute:!0}):null,!this.checkable&&this.mergedBordered?e(`div`,{class:`${n}-tag__border`,style:{borderColor:a}}):null)}}),W={key:0,class:`ml-4px`},G={key:0,class:`mt-2px text-12px opacity-80`},K=F({name:`ProjectDataFreshness`,__name:`data-freshness`,props:{dataAt:{default:null},delayedAfterMinutes:{default:10},expiredAfterMinutes:{default:30},compact:{type:Boolean,default:!1}},setup(e){let r=e,i=T(Date.now()),c,l=C(()=>{if(!r.dataAt)return null;let e=new Date(r.dataAt).getTime();return Number.isFinite(e)?e:null}),f=C(()=>{if(l.value===null)return{label:`数据时间未知`,type:`warning`,detail:`中央接口没有提供有效的同步时间`};let e=i.value-l.value;if(e<-300*1e3)return{label:`数据时间未知`,type:`warning`,detail:`同步时间晚于当前时间，可能存在时钟偏差`};let t=Math.max(r.delayedAfterMinutes,1)*60*1e3,n=Math.max(r.expiredAfterMinutes,r.delayedAfterMinutes)*60*1e3;return e>n?{label:`数据过期`,type:`error`,detail:`已超过 ${Math.round(n/6e4)} 分钟未同步`}:e>t?{label:`数据延迟`,type:`warning`,detail:`已超过 ${Math.round(t/6e4)} 分钟未同步`}:{label:`数据新鲜`,type:`success`,detail:`最近一次同步在可接受范围内`}}),p=C(()=>l.value===null?``:new Date(l.value).toLocaleString(`zh-CN`,{hour12:!1}));return A(()=>{r.compact||(c=setInterval(()=>{i.value=Date.now()},6e4))}),n(()=>{c&&clearInterval(c)}),(n,r)=>{let i=U,c=w;return h(),s(c,null,{trigger:j(()=>[d(i,{size:`small`,type:f.value.type},{default:j(()=>[a(x(f.value.label)+` `,1),!e.compact&&p.value?(h(),t(`span`,W,`· `+x(p.value),1)):u(``,!0)]),_:1},8,[`type`])]),default:j(()=>[o(`div`,null,x(f.value.detail),1),p.value?(h(),t(`div`,G,`数据时间：`+x(p.value),1)):u(``,!0)]),_:1})}}});export{U as n,K as t};